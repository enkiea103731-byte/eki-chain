package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
)

var balances = make(map[string]*big.Int)

func main() {
	http.HandleFunc("/", handleRPC)
	http.HandleFunc("/mint", handleMint)
	
	port := "10000"
	if p := getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("EKI Chain 70101 rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getenv(k string) string {
	// pega PORT do Render
	return strings.TrimSpace(func() string {
		b, _ := json.Marshal(map[string]string{})
		_ = b
		return ""
	}())
}

func handleMint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	var req struct {
		To     string `json:"to"`
		Amount string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"json invalido"}`, 400)
		return
	}
	amount := new(big.Int)
	amount.SetString(req.Amount, 10)
	
	to := strings.ToLower(req.To)
	balances[to] = amount
	
	resp := map[string]interface{}{
		"status":  "minted",
		"to":      req.To,
		"balance": req.Amount,
		"chainId": "70101",
	}
	json.NewEncoder(w).Encode(resp)
	fmt.Printf("MINT %s -> %s\n", req.To, req.Amount)
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method == "GET" {
		// Info da chain que apareceu na sua foto
		info := map[string]interface{}{
			"chain": "EKI", "chainId": 70101, "chainSlug": "eki",
			"custody": map[string]string{
				"BNB": "0x03664899007a53F086D545c34063E9CBfb29C9C2",
				"ETH": "0x4DcC1EA8ae08f1bdaF5660ab37Adc2c73c93cf",
				"BTC": "bc1qe03sft8p9aqrsvfn8qg8sm7afh3g30ru5zvt",
			},
			"feeWallet": "EKI_31e3da3a96a0af3",
			"feeRate": 0.003,
			"rpc": []string{"https://eki-chain.onrender.com"},
			"nativeCurrency": map[string]interface{}{"name":"EKI","symbol":"EKI","decimals":18},
		}
		json.NewEncoder(w).Encode(info)
		return
	}

	var req struct {
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
		ID     interface{}   `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var result interface{}
	switch req.Method {
	case "eth_chainId":
		result = "0x11165" // 70101
	case "eth_getBalance":
		addr := ""
		if len(req.Params) > 0 {
			if s, ok := req.Params[0].(string); ok {
				addr = strings.ToLower(s)
			}
		}
		bal, ok := balances[addr]
		if !ok {
			result = "0x0"
		} else {
			result = "0x" + bal.Text(16)
		}
	case "eth_blockNumber", "eth_gasPrice":
		result = "0x1"
	case "net_version":
		result = "70101"
	default:
		result = "0x0"
	}

	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      req.ID,
		"result":  result,
	}
	json.NewEncoder(w).Encode(resp)
}
