package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

var balances = make(map[string]*big.Int)

func main() {
	http.HandleFunc("/mint", handleMint)
	http.HandleFunc("/", handleRPC)
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Println("EKI 70101 rodando na porta", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleMint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
	var req struct{ To string `json:"to"`; Amount string `json:"amount"` }
	json.NewDecoder(r.Body).Decode(&req)
	amt := new(big.Int)
	amt.SetString(req.Amount, 10)
	balances[strings.ToLower(req.To)] = amt
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "minted", "to": req.To, "balance": req.Amount, "chainId": "70101",
	})
	fmt.Println("MINT", req.To, req.Amount)
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" { w.Write([]byte("{}")); return }
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"chain":"EKI","chainId":70101,
			"custody": map[string]string{"BNB":"0x03664899007a53F086D545c34063E9CBfb29C9C2"},
			"feeWallet":"EKI_31e3da3a96a0af3","nativeCurrency":map[string]interface{}{"name":"EKI","symbol":"EKI","decimals":18},
		})
		return
	}
	var req struct{ Method string `json:"method"`; Params []interface{} `json:"params"`; ID interface{} `json:"id"` }
	json.NewDecoder(r.Body).Decode(&req)
	var result interface{}
	switch req.Method {
	case "eth_chainId": result = "0x11165"
	case "net_version": result = "70101"
	case "eth_getBalance":
		addr := ""
		if len(req.Params)>0 { if s,ok:=req.Params[0].(string); ok { addr=strings.ToLower(s) } }
		if b, ok:=balances[addr]; ok { result="0x"+b.Text(16) } else { result="0x0" }
	default: result="0x0"
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc":"2.0","id":req.ID,"result":result})
}
