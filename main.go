package main
import (
 "encoding/json"
 "log"
 "net/http"
 "strings"
)
const CHAIN_ID = 701
const FEE_WALLET = "EKI_31e3da3ae96a0af3"
var custody = map[string]string{
 "BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
 "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
 "BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
 "SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm",
}
func applyFee(a float64) (float64,float64){fee:=a*0.003;return a-fee,fee}
func main(){
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{"address":FEE_WALLET,"chainId":CHAIN_ID,"name":"EKIChain","status":"Online","rpc":"https://eki-chain.onrender.com","symbol":"EKI"})
 })
 http.HandleFunc("/bridge/deposit", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  chain:=strings.ToUpper(r.URL.Query().Get("chain"))
  toUser,fee:=applyFee(1.0)
  if chain==""{
   json.NewEncoder(w).Encode(map[string]interface{}{"feeRate":0.003,"feeWallet":FEE_WALLET,"custody":custody})
   return
  }
  addr,ok:=custody[chain]
  if !ok{
   w.WriteHeader(400)
   json.NewEncoder(w).Encode(map[string]string{"error":"use BTC, ETH, BNB, SOL"})
   return
  }
  json.NewEncoder(w).Encode(map[string]interface{}{"chain":chain,"custody":addr,"deposit":1.0,"toUser":toUser,"fee":fee,"feeWallet":FEE_WALLET})
 })
 log.Println("EKI Chain 701 rodando na 8080 com 4 cadeias")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
