package main
import (
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "os"
 "strconv"
 "sync/atomic"
)
var blockNumber int64 = 1
var totalSupplyHex = "0x33B2E3C9FD0803CE8000000"
func main(){
 port:=os.Getenv("PORT")
 if port==""{port="10000"}
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  json.NewEncoder(w).Encode(map[string]interface{}{"chainId":70101,"name":"EKI 70101"})
 })
 http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS")
  w.Header().Set("Access-Control-Allow-Headers","Content-Type")
  if r.Method=="OPTIONS"{return}
  var req map[string]interface{}
  json.NewDecoder(r.Body).Decode(&req)
  method,_:=req["method"].(string)
  id:=req["id"]
  atomic.AddInt64(&blockNumber,1)
  bn:="0x"+strconv.FormatInt(atomic.LoadInt64(&blockNumber),16)
  var result interface{}
  switch method{
  case "eth_chainId": result="0x11165"
  case "net_version": result="70101"
  case "eth_blockNumber": result=bn
  case "eth_getBalance": result=totalSupplyHex
  case "eth_getTransactionCount": result="0x1"
  case "eth_gasPrice": result="0x3B9ACA00"
  case "eth_estimateGas": result="0x5208"
  default: result="0x0"
  }
  json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc":"2.0","id":id,"result":result})
 })
 fmt.Printf("EKI 70101 Online on %s - FIX 1Q pra todas carteiras\n",port)
 log.Fatal(http.ListenAndServe(":"+port,nil))
}
