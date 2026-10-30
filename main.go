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
	// 1 Quatrilhão = 10^15 tokens * 10^18 decimais = 10^33 wei
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
	log.Printf("EKI MAINNET 70101 - 1 QUATRILHAO - porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { handleRPC(w, r); return }
	setCORS(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"network": "EKI MAINNET",
		"chainId": 70101,
		"chainIdHex": "0x11165",
		"rpc": "https://eki-chain.onrender.com/rpc",
		"symbol": "EKI",
		"decimals": 18,
		"totalSupply": "1000000000000000",
		"totalSupplyWei": ONE_QUAD.String(),
		"custodyBNB": MY_ADDR,
		"feeWallet": "EKI_31e3da3a96a0af3",
		"endpoints": []string{"/rpc", "/mint", "/transfer", "/balance?address=0x...", "/price", "/explorer"},
	})
}

func handlePrice(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pair": "EKI/BNB",
		"price_bnb": "0.00000001",
		"lowest_wei": "10000000",
		"marketCap": "1 Quatrilhao EKI",
		"buy_example": "Envie BNB para " + MY_ADDR + " e receba EKI",
	})
}

func handleMint(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
	var req struct { To string `json:"to"`; Amount string `json:"amount"` }
	json.NewDecoder(r.Body).Decode(&req)
	if req.To == "" { req.To = MY_ADDR }
	amt := new(big.Int)
	if req.Amount == "" { amt.Set(ONE_QUAD) } else { amt.SetString(req.Amount, 10) }
	mu.Lock()
	balances[strings.ToLower(req.To)] = new(big.Int).Set(amt)
	mu.Unlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"status":"MAINNET_MINTED","to":req.To,"amount":amt.String(),"human":"1 Quatrilhao","chainId":70101})
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	var req struct { From, To, Amount string `json:"from"; json:"to"; json:"amount"` }
	json.NewDecoder(r.Body).Decode(&req)
	from := strings.ToLower(req.From); to := strings.ToLower(req.To)
	amt := new(big.Int); amt.SetString(req.Amount, 10)
	mu.Lock()
	defer mu.Unlock()
	bal, ok := balances[from]
	if !ok || bal.Cmp(amt) < 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"error":"saldo insuficiente"}); return
	}
	balances[from] = new(big.Int).Sub(bal, amt)
	if balances[to] == nil { balances[to] = big.NewInt(0) }
	balances[to] = new(big.Int).Add(balances[to], amt)
	json.NewEncoder(w).Encode(map[string]interface{}{"status":"sent","from":req.From,"to":req.To,"amount":req.Amount})
}

func handleBalance(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	addr := strings.ToLower(r.URL.Query().Get("address"))
	if addr == "" { addr = strings.ToLower(MY_ADDR) }
	mu.RLock()
	bal := balances[addr]
	mu.RUnlock()
	if bal == nil { bal = big.NewInt(0) }
	json.NewEncoder(w).Encode(map[string]interface{}{"address":addr,"balanceWei":bal.String(),"balanceEKI":"1000000000000000","chainId":70101})
}

func handleExplorer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body style="font-family:sans-serif;background:#111;color:#0f0;padding:20px"><h1>EKI MAINNET 70101</h1><h2>1 QUATRILHAO MINTADO</h2><p>RPC: https://eki-chain.onrender.com/rpc</p><p>ChainId: 70101 (0x11165)</p><p>Supply: %s</p><p>Titular: %s</p></body></html>`, ONE_QUAD.String(), MY_ADDR)
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
	if r.Method == "GET" { handleInfo(w, r); return }
	var req struct { Method string `json:"method"`; Params []interface{} `json:"params"`; ID interface{} `json:"id"` }
	json.NewDecoder(r.Body).Decode(&req)
	var result interface{}
	switch req.Method {
	case "eth_chainId": result = "0x11165"
	case "net_version": result = "70101"
	case "eth_gasPrice": result = "0x0"
	case "eth_blockNumber": result = "0x1"
	case "eth_getBalance":
		addr := ""
		if len(req.Params) > 0 { if s, ok := req.Params[0].(string); ok { addr = strings.ToLower(s) } }
		mu.RLock()
		b := balances[addr]
		mu.RUnlock()
		if b != nil { result = "0x" + b.Text(16) } else { result = "0x0" }
	case "eth_sendTransaction", "eth_call": result = "0x0"
	default: result = "0x0"
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc":"2.0","id":req.ID,"result":result})
}
