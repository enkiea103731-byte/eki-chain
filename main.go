package main

import (
"encoding/json"
"fmt"
"log"
"math/big"
"net/http"
"os"
"strings"
"sync"
)

var (
balances = make(map[string]*big.Int)
mu       sync.RWMutex
ONE_QUAD = new(big.Int).Exp(big.NewInt(10), big.NewInt(33), nil)
CUSTODY = map[string]string{
"BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
"ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
"BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
"SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm",
"EKI_MAIN": "0x03664899007a53f086d545c34063e9cbfb29c9c2",
}
)

func init() {
for _, addr := range CUSTODY {
balances[strings.ToLower(addr)] = new(big.Int).Set(ONE_QUAD)
}
balances[strings.ToLower(CUSTODY["EKI_MAIN"])] = new(big.Int).Set(ONE_QUAD)
}

func main() {
http.HandleFunc("/", handleRoot)
http.HandleFunc("/price", handlePrice)
http.HandleFunc("/balance", handleBalance)
http.HandleFunc("/mint", handleMint)
http.HandleFunc("/explorer", handleExplorer)
http.HandleFunc("/rpc", handleRPC)
port := os.Getenv("PORT")
if port == "" { port = "10000" }
log.Printf("EKI MAINNET 4 CUSTODIAS ONLINE %s", port)
log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setCORS(w http.ResponseWriter) {
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
w.Header().Set("Content-Type", "application/json")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
if r.Method == "POST" { handleRPC(w, r); return }
setCORS(w)
json.NewEncoder(w).Encode(map[string]interface{}{
"network": "EKI MAINNET", "chainId": 70101, "chainIdHex": "0x11165",
"totalSupply": "1000000000000000", "totalSupplyWei": ONE_QUAD.String(),
"custody": CUSTODY,
"rpc": "https://eki-chain.onrender.com/rpc",
"status": "MAINNET ONLINE - 4 MOEDAS ATIVAS",
})
}

func handlePrice(w http.ResponseWriter, r *http.Request) {
setCORS(w)
json.NewEncoder(w).Encode(map[string]interface{}{
"EKI/BNB": "0.00000001", "supply": "1 Quatrilhao", "custody": CUSTODY,
})
}

func handleBalance(w http.ResponseWriter, r *http.Request) {
setCORS(w)
addr := strings.ToLower(r.URL.Query().Get("address"))
if addr == "" { addr = strings.ToLower(CUSTODY["EKI_MAIN"]) }
mu.RLock(); b := balances[addr]; mu.RUnlock()
if b == nil { b = big.NewInt(0) }
json.NewEncoder(w).Encode(map[string]interface{}{"address": addr, "balanceWei": b.String(), "balanceEKI": "1000000000000000", "custody": CUSTODY})
}

func handleMint(w http.ResponseWriter, r *http.Request) {
setCORS(w)
if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
var req struct{ To string `json:"to"` }
json.NewDecoder(r.Body).Decode(&req)
if req.To == "" { req.To = CUSTODY["EKI_MAIN"] }
mu.Lock(); balances[strings.ToLower(req.To)] = new(big.Int).Set(ONE_QUAD); mu.Unlock()
json.NewEncoder(w).Encode(map[string]interface{}{"status": "MINTED", "to": req.To, "supply": ONE_QUAD.String()})
}

func handleExplorer(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
fmt.Fprintf(w, `<html><body style="background:#111;color:#0f0;padding:20px;font-family:sans-serif"><h1>EKI MAINNET 70101 - ONLINE</h1><h2>1 QUATRILHAO</h2><p>BTC: %s</p><p>ETH: %s</p><p>BNB: %s</p><p>SOL: %s</p><p>EKI: %s</p></body></html>`, CUSTODY["BTC"], CUSTODY["ETH"], CUSTODY["BNB"], CUSTODY["SOL"], CUSTODY["EKI_MAIN"])
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
setCORS(w)
if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
if r.Method == "GET" { handleRoot(w, r); return }
var req struct{ Method string `json:"method"`; Params []interface{} `json:"params"`; ID interface{} `json:"id"` }
json.NewDecoder(r.Body).Decode(&req)
var result interface{}
switch req.Method {
case "eth_chainId": result = "0x11165"
case "net_version": result = "70101"
case "eth_getBalance":
addr := ""; if len(req.Params) > 0 { if s, ok := req.Params[0].(string); ok { addr = strings.ToLower(s) } }
mu.RLock(); b := balances[addr]; mu.RUnlock()
if b != nil { result = "0x" + b.Text(16) } else { result = "0x0" }
default: result = "0x0"
}
json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc": "2.0", "id": req.ID, "result": result})
}
