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
MY_ADDR  = "0x03664899007a53f086d545c34063e9cbfb29c9c2"
)

func init() {
balances[strings.ToLower(MY_ADDR)] = new(big.Int).Set(ONE_QUAD)
}

func main() {
http.HandleFunc("/", handleInfo)
http.HandleFunc("/mint", handleMint)
http.HandleFunc("/transfer", handleTransfer)
http.HandleFunc("/balance", handleBalance)
http.HandleFunc("/price", handlePrice)
http.HandleFunc("/rpc", handleRPC)
http.HandleFunc("/explorer", handleExplorer)
port := os.Getenv("PORT")
if port == "" { port = "10000" }
log.Printf("EKI MAINNET 1 QUATRILHAO na porta %s", port)
log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setCORS(w http.ResponseWriter) {
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
w.Header().Set("Content-Type", "application/json")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
if r.URL.Path != "/" && r.URL.Path != "/rpc" {
w.WriteHeader(404); return
}
if r.URL.Path == "/rpc" && r.Method == "GET" {
setCORS(w)
json.NewEncoder(w).Encode(map[string]interface{}{"chain":"EKI","chainId":70101,"status":"MAINNET ONLINE"})
return
}
if strings.Contains(r.Header.Get("Content-Type"), "application/json") || r.Method == "POST" {
handleRPC(w, r); return
}
setCORS(w)
json.NewEncoder(w).Encode(map[string]interface{}{
"network": "EKI MAINNET",
"chainId": 70101,
"chainIdHex": "0x11165",
"rpc": "https://eki-chain.onrender.com/rpc",
"totalSupply": "1000000000000000",
"totalSupplyWei": ONE_QUAD.String(),
})
}

func handlePrice(w http.ResponseWriter, r *http.Request) { setCORS(w); json.NewEncoder(w).Encode(map[string]interface{}{"pair":"EKI/BNB","price_bnb":"0.00000001","supply":"1 Quatrilhao"}) }
func handleMint(w http.ResponseWriter, r *http.Request) {
setCORS(w)
if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
var req struct{ To string `json:"to"`; Amount string `json:"amount"` }
json.NewDecoder(r.Body).Decode(&req)
if req.To == "" { req.To = MY_ADDR }
amt := new(big.Int)
if req.Amount == "" { amt.Set(ONE_QUAD) } else { amt.SetString(req.Amount, 10) }
mu.Lock(); balances[strings.ToLower(req.To)] = new(big.Int).Set(amt); mu.Unlock()
json.NewEncoder(w).Encode(map[string]interface{}{"status":"MAINNET_MINTED","to":req.To,"amount":amt.String()})
}
func handleTransfer(w http.ResponseWriter, r *http.Request) {
setCORS(w)
var req struct{ From string `json:"from"`; To string `json:"to"`; Amount string `json:"amount"` }
json.NewDecoder(r.Body).Decode(&req)
mu.Lock(); defer mu.Unlock()
from := strings.ToLower(req.From); to := strings.ToLower(req.To)
amt := new(big.Int); amt.SetString(req.Amount, 10)
if balances[from] == nil || balances[from].Cmp(amt) < 0 { json.NewEncoder(w).Encode(map[string]interface{}{"error":"saldo insuficiente"}); return }
balances[from] = new(big.Int).Sub(balances[from], amt)
if balances[to] == nil { balances[to] = big.NewInt(0) }
balances[to] = new(big.Int).Add(balances[to], amt)
json.NewEncoder(w).Encode(map[string]interface{}{"status":"sent"})
}
func handleBalance(w http.ResponseWriter, r *http.Request) {
setCORS(w)
addr := strings.ToLower(r.URL.Query().Get("address"))
if addr == "" { addr = strings.ToLower(MY_ADDR) }
mu.RLock(); b := balances[addr]; mu.RUnlock()
if b == nil { b = big.NewInt(0) }
json.NewEncoder(w).Encode(map[string]interface{}{"address":addr,"balanceWei":b.String(),"balanceEKI":"1000000000000000"})
}
func handleExplorer(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
fmt.Fprintf(w, `<h1>EKI MAINNET ONLINE</h1><p>1 QUATRILHAO</p><p>%s</p>`, ONE_QUAD.String())
}
func handleRPC(w http.ResponseWriter, r *http.Request) {
setCORS(w)
if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
if r.Method == "GET" { handleInfo(w, r); return }
var req struct{ Method string `json:"method"`; Params []interface{} `json:"params"`; ID interface{} `json:"id"` }
json.NewDecoder(r.Body).Decode(&req)
var result interface{}
switch req.Method {
case "eth_chainId": result = "0x11165"
case "net_version": result = "70101"
case "eth_getBalance":
addr := ""; if len(req.Params)>0 { if s,ok:=req.Params[0].(string); ok { addr=strings.ToLower(s) } }
mu.RLock(); b:=balances[addr]; mu.RUnlock()
if b!=nil { result="0x"+b.Text(16) } else { result="0x0" }
default: result="0x0"
}
json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc":"2.0","id":req.ID,"result":result})
}
