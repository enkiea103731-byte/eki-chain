package main

import (
  "encoding/json"
  "net/http"
  "os"
)

func main() {
  chainInfo := map[string]interface{}{
    "chain": "EKI", "chainId": 70101, "chainSlug": "eki",
    "custody": map[string]string{
      "BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
      "BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
      "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
      "SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm",
    },
    "explorers": []interface{}{}, "faucets": []interface{}{},
    "feeRate": 0.003, "feeWallet": "EKI_31e3da3ae96a0af3",
    "infoURL": "https://eki-chain.onrender.com",
    "name": "EKI Chain",
    "nativeCurrency": map[string]interface{}{"decimals": 18, "name": "EKI", "symbol": "EKI"},
    "networkId": 70101,
    "rpc": []string{"https://eki-chain.onrender.com"},
    "shortName": "eki", "status": "Online",
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      var req map[string]interface{}
      json.NewDecoder(r.Body).Decode(&req)
      id := req["id"]
      method, _ := req["method"].(string)
      result := "0x0"
      if method == "eth_chainId" { result = "0x11235" }
      if method == "net_version" { result = "70101" }
      json.NewEncoder(w).Encode(map[string]interface{}{
        "jsonrpc": "2.0", "id": id, "result": result,
      })
      return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chainInfo)
  })

  port := os.Getenv("PORT")
  if port == "" { port = "3000" }
  http.ListenAndServe(":"+port, nil)
}
