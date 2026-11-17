package main

import (
  "encoding/json"
"math/big"
  "net/http"
  "strings"
)

func rpcHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  if r.Method == "OPTIONS" { w.WriteHeader(200); return }
  var req struct {
    ID any `json:"id"`
    Method string `json:"method"`
    Params []json.RawMessage `json:"params"`
  }
  json.NewDecoder(r.Body).Decode(&req)
  code := "0x608060405234801561001057600080fd5b50"
  isContract := false
  if len(req.Params) > 0 {
    a := strings.ToLower(strings.Trim(string(req.Params[0]), "\""))
    if a == "0x227a5639865eef7c053ac937d169448b4178bf52" || a == "0x9096ff1c3fc2487396fa88bf9fd4085dc2f38619" || a == "0xba218e9c339d43599ed8b4d56c0ab406a7d33e79" || a == "0xc5018b51be6e2d72e5758a5dcaaaf0027fdd7e65" {
      isContract = true
    }
  }
  var result any
  switch req.Method {
  case "eth_chainId":
    result = "0x111d5"
  case "net_version":
    result = "70101"
  case "eth_blockNumber":
    result = "0x3e8"
  case "eth_gasPrice":
    result = "0x3b9aca00"
  case "eth_getBalance":
  // LIBERADO - saldo real
  mu.Lock()
  bal := 0.0
  addr := ""
  if len(req.Params) > 0 { var s string; json.Unmarshal(req.Params[0], &s); addr = strings.ToLower(s) }
  if addr == "0x50b1efba9ad83f294fd4157a931c050109ffe30a" { bal = c.Saldos["enki71"]; if bal == 0 { bal = c.SaldoEnki } } else if addr == "0xb8ce0fe09ac0cd186a1006a4925c5af644766b45" { bal = c.SaldoEnki } else if addr == "0x0909ff6f1c3f248736f88b9ff64d4085dc2f38619" { bal = c.Saldos["Marta"] }
  mu.Unlock()
  wei := big.NewInt(int64(bal))
  wei.Mul(wei, big.NewInt(1000000000000000000))
  result = "0x" + wei.Text(16)
  case "eth_getTransactionCount":
    result = "0x0"
  case "eth_getCode":
    if isContract { result = code } else { result = "0x" }
  case "eth_call":
    data := ""
    if len(req.Params) > 0 { data = strings.ToLower(string(req.Params[0])) }
    if strings.Contains(data, "06fdde03") {
      result = "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003454b490000000000000000000000000000000000000000000000000000000000"
    } else if strings.Contains(data, "95d89b41") {
      result = "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003454b490000000000000000000000000000000000000000000000000000000000"
    } else if strings.Contains(data, "313ce567") {
      result = "0x0000000000000000000000000000000000000000000000000000000000000012"
    } else if strings.Contains(data, "18160ddd") {
      result = "0x00000000000000000000000000000000000000000000000000d3c21bcecceda1000000"
    } else if strings.Contains(data, "70a08231") {
      result = "0x00000000000000000000000000000000000000000000000000d3c21bcecceda1000000"
    } else {
      result = "0x"
    }
  case "eth_estimateGas":
    result = "0x5208"
  case "eth_getBlockByNumber":
    result = map[string]any{"number": "0x3e8", "hash": "0x111d5000000000000000000000000000000000000000000000000000000000000001", "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000", "timestamp": "0x692a0000", "transactions": []any{}, "gasLimit": "0x1c9c380", "gasUsed": "0x0", "miner": "0x0000000000000000000000000000000000000000", "difficulty": "0x1", "totalDifficulty": "0x3e8", "size": "0x200", "uncles": []any{}, "nonce": "0x0000000000000000", "extraData": "0x"}
  case "eth_getBlockByHash":
    result = map[string]any{"number": "0x3e8", "hash": "0x111d5000000000000000000000000000000000000000000000000000000000000001", "timestamp": "0x692a0000", "transactions": []any{}}
  case "eth_getTransactionByHash":
    result = nil
  case "eth_getTransactionReceipt":
    result = nil
  case "eth_getLogs":
    result = []any{}
  case "eth_syncing":
    result = false
  case "eth_accounts":
    result = []any{}
  case "net_listening":
    result = true
  case "web3_clientVersion":
    result = "EKI/1.0"
  default:
    result = "0x0"
  }
  json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
}
