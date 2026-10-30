package main

import (
 "encoding/json"
 "io"
 "log"
 "net/http"
 "strings"
)

const CHAIN_ID = 701
const CHAIN_ID_HEX = "0x2bd"
const FEE_WALLET = "EKI_31e3da3ae96a0af3"

var custody = map[string]string{
 "BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
 "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
 "BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
 "SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm",
}

type RPCReq struct {
 ID     interface{} `json:"id"`
 Method string      `json:"method"`
 Params json.RawMessage `json:"params"`
}

func writeRPC(w http.ResponseWriter, id interface{}, result interface{}) {
 w.Header().Set("Content-Type","application/json")
 w.Header().Set("Access-Control-Allow-Origin","*")
 json.NewEncoder(w).Encode(map[string]interface{}{
  "jsonrpc":"2.0","id":id,"result":result,
 })
}

func main(){
 // CORS
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Headers","Content-Type")
  w.Header().Set("Access-Control-Allow-Methods","GET,POST,OPTIONS")
  if r.Method=="OPTIONS"{ w.WriteHeader(200); return }

  // CHAINLIST / Metamask vai fazer POST com eth_chainId
  if r.Method=="POST"{
   body,_:=io.ReadAll(r.Body)
   var req RPCReq
   json.Unmarshal(body,&req)
   switch req.Method {
   case "eth_chainId":
    writeRPC(w,req.ID,CHAIN_ID_HEX)
    return
   case "net_version":
    writeRPC(w,req.ID,"701")
    return
   case "eth_blockNumber":
    writeRPC(w,req.ID,"0x1")
    return
   case "eth_gasPrice":
    writeRPC(w,req.ID,"0x0")
    return
   case "eth_getBlockByNumber","eth_getBlockByHash":
    writeRPC(w,req.ID,map[string]interface{}{"number":"0x1","hash":"0x0","parentHash":"0x0"})
    return
   case "eth_accounts","eth_requestAccounts":
    writeRPC(w,req.ID,[]string{})
    return
   case "eth_getBalance":
    writeRPC(w,req.ID,"0x0")
    return
   default:
    writeRPC(w,req.ID,"0x0")
    return
   }
  }

  // GET / -> Metadata Chainlist + seu custody (os 4)
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "name":"EKI Chain",
   "chain":"EKI",
   "chainId":CHAIN_ID,
   "networkId":CHAIN_ID,
   "shortName":"eki",
   "chainSlug":"eki",
   "nativeCurrency":map[string]interface{}{"name":"EKI","symbol":"EKI","decimals":18},
   "rpc":[]string{"https://eki-chain.onrender.com"},
   "faucets":[]string{},
   "infoURL":"https://eki-chain.onrender.com",
   "explorers":[]string{},
   "status":"Online",
   "feeRate":0.003,
   "feeWallet":FEE_WALLET,
   "custody":custody,
  })
 })

 http.HandleFunc("/bridge/deposit", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  w.Header().Set("Access-Control-Allow-Origin","*")
  chain:=strings.ToUpper(r.URL.Query().Get("chain"))
  if chain==""{
   json.NewEncoder(w).Encode(map[string]interface{}{"custody":custody,"feeRate":0.003,"feeWallet":FEE_WALLET})
   return
  }
  json.NewEncoder(w).Encode(map[string]interface{}{"chain":chain,"custody":custody[chain],"fee":0.003,"feeWallet":FEE_WALLET,"toUser":0.997})
 })

 log.Println("EKI Chain 701 CHAINLIST-READY com 4 cadeias ativo na 8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
