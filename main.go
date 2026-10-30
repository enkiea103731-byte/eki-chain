package main

import (
"encoding/json"
"log"
"math/big"
"net/http"
"os"
"strings"
"sync"
)

var (
balances = make(map[string]*big.Int)
mu sync.RWMutex
ONE_QUAD = new(big.Int).Exp(big.NewInt(10), big.NewInt(33), nil)
CUSTODY = map[string]string{
"BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
"ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
"BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
"SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm",
"EKI": "0x03664899007a53f086d545c34063e9cbfb29c9c2",
}
)

func init(){ for _,a:=range CUSTODY{ balances[strings.ToLower(a)]=new(big.Int).Set(ONE_QUAD) } }

func main(){
http.HandleFunc("/", handleRoot)
http.HandleFunc("/rpc", handleRPC)
http.HandleFunc("/chainlist.json", handleChainlist)
http.HandleFunc("/chain.json", handleChainlist)
http.HandleFunc("/explorer", func(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html")
w.Write([]byte("<h1>EKI MAINNET 70101 - 1 QUATRILHAO ONLINE</h1><p>RPC: https://eki-chain.onrender.com/rpc</p>"))
})
port:=os.Getenv("PORT"); if port==""{port="10000"}
log.Fatal(http.ListenAndServe(":"+port,nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request){
if r.Method=="POST"{ handleRPC(w,r); return }
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"chainId":70101,"network":"EKI MAINNET","custody":CUSTODY,"totalSupplyWei":ONE_QUAD.String(),"chainlistFile":"https://eki-chain.onrender.com/chainlist.json"})
}

func handleChainlist(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
"name": "EKI Mainnet",
"chain": "EKI",
"icon": "eki",
"rpc": []string{"https://eki-chain.onrender.com/rpc"},
"faucets": []string{},
"nativeCurrency": map[string]interface{}{"name":"EKI","symbol":"EKI","decimals":18},
"infoURL": "https://eki-chain.onrender.com",
"shortName": "eki",
"chainId": 70101,
"networkId": 70101,
"slip44": 1,
"explorers": []map[string]string{{"name":"EKI Explorer","url":"https://eki-chain.onrender.com/explorer","standard":"none"}},
})
}

func handleRPC(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Access-Control-Allow-Headers","Content-Type")
w.Header().Set("Access-Control-Allow-Methods","GET, POST, OPTIONS")
w.Header().Set("Content-Type","application/json")
if r.Method=="OPTIONS"{ w.Write([]byte("{}")); return }
if r.Method=="GET"{ handleRoot(w,r); return }
var req struct{ Method string `json:"method"`; Params []interface{} `json:"params"`; ID interface{} `json:"id"` }
json.NewDecoder(r.Body).Decode(&req)
var result interface{}
switch req.Method {
case "eth_chainId": result="0x11165"
case "net_version": result="70101"
case "eth_blockNumber": result="0x1"
case "eth_gasPrice": result="0x3B9ACA00"
case "eth_estimateGas": result="0x5208"
case "eth_getTransactionCount": result="0x1"
case "eth_getBalance":
addr:=""; if len(req.Params)>0{ if s,ok:=req.Params[0].(string); ok{ addr=strings.ToLower(s) } }
mu.RLock(); b:=balances[addr]; mu.RUnlock()
if b!=nil{ result="0x"+b.Text(16) } else { result="0x0" }
default: result="0x0"
}
json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc":"2.0","id":req.ID,"result":result})
}
