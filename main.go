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
  json.NewEncoder(w).Encode(map[string]interface{}{
   "chainId":CHAIN_ID,"feeWallet":FEE_WALLET,"custody":custody,"feeRate":0.003,"status":"Online",
  })
 })
 http.HandleFunc("/bridge/deposit", func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  chain:=strings.ToUpper(r.URL.Query().Get("chain"))
  toUser,fee:=applyFee(1.0)
  if chain==""{
   json.NewEncoder(w).Encode(map[string]interface{}{"custody":custody,"feeRate":0.003,"feeWallet":FEE_WALLET,"example":"?chain=BTC | ETH | BNB | SOL"})
   return
  }
  addr:=custody[chain]
  json.NewEncoder(w).Encode(map[string]interface{}{"chain":chain,"custody":addr,"toUser":toUser,"fee":fee,"feeWallet":FEE_WALLET})
 })
 log.Println("EKI Chain 701 com 4 cadeias ativo")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
