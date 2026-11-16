package main

import (
  "encoding/json"
  "net/http"
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
  block := map[string]any{"number":"0x1f","hash":"0x111d5000000000000000000000000000000000000000000000000000000000001","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","timestamp":"0x692a0000","transactions":[]any{},"gasLimit":"0x1c9c380","gasUsed":"0x0","miner":"0x0000000000000000000000000000000000000000","difficulty":"0x1","totalDifficulty":"0x1f","size":"0x200","uncles":[]any{},"nonce":"0x0000000000000000","sha3Uncles":"0x0000000000000000000000000000000000000000000000000000000000000000","logsBloom":"0x00","transactionsRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","stateRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","receiptsRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","extraData":"0x"}
  var result any
  switch req.Method {
  case "eth_chainId": result = "0x111d5"
  case "net_version": result = "70101"
  case "eth_blockNumber": result = "0x1f"
  case "eth_gasPrice": result = "0x3b9aca00"
  case "eth_maxPriorityFeePerGas": result = "0x3b9aca00"
  case "eth_feeHistory": result = map[string]any{"baseFeePerGas":[]any{"0x3b9aca00"},"gasUsedRatio":[]any{0},"reward":[]any{[]any{"0x0"}},"oldestBlock":"0x1f"}
  case "eth_getBalance": result = "0x0"
  case "eth_getCode": result = "0x"
  case "eth_getTransactionCount": result = "0x0"
  case "eth_getStorageAt": result = "0x0000000000000000000000000000000000000000000000000000000000000000"
  case "eth_call": result = "0x"
  case "eth_estimateGas": result = "0x5208"
  case "eth_getBlockByNumber": result = block
  case "eth_getBlockByHash": result = block
  case "eth_getTransactionByHash": result = nil
  case "eth_getTransactionReceipt": result = nil
  case "eth_getLogs": result = []any{}
  case "eth_syncing": result = false
  case "eth_accounts": result = []any{}
  case "net_listening": result = true
  case "web3_clientVersion": result = "EKI/1.0"
  case "web3_sha3": result = "0x0000000000000000000000000000000000000000000000000000000000000000"
  default: result = "0x0"
  }
  json.NewEncoder(w).Encode(map[string]any{"jsonrpc":"2.0","id":req.ID,"result":result})
}
