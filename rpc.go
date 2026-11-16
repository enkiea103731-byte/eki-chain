package main

import (
  "encoding/json"
  "net/http"
)

func rpcHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == "OPTIONS" { w.WriteHeader(200); return }
  var req struct {
    ID any `json:"id"`
    Method string `json:"method"`
    Params []json.RawMessage `json:"params"`
  }
  json.NewDecoder(r.Body).Decode(&req)
  var result any
  switch req.Method {
  case "eth_chainId":
    result = "0x111d5"
  case "net_version":
    result = "70101"
  case "eth_blockNumber":
    result = "0x0"
  case "eth_gasPrice":
    result = "0x0"
  case "eth_getBalance":
    result = "0x0"
  case "eth_getCode":
    result = "0x"
  case "eth_getTransactionCount":
    result = "0x0"
  case "eth_call":
    result = "0x"
  case "eth_estimateGas":
    result = "0x5208"
  case "eth_getBlockByNumber":
    result = map[string]any{"number":"0x0","hash":"0x0000000000000000000000000000000000000000000000000000000000000000","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","timestamp":"0x0","transactions":[]any{}}
  case "web3_clientVersion":
    result = "EKI/1.0"
  default:
    result = "0x0"
  }
  json.NewEncoder(w).Encode(map[string]any{"jsonrpc":"2.0","id":req.ID,"result":result})
}
