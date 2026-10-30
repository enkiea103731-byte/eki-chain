package main

import (
"encoding/json"
"log"
"net/http"
)

const CHAIN_ID = 701
const FEE_WALLET = "EKI_31e3da3ae96a0af3"

const (
CUSTODY_BTC = "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt"
CUSTODY_ETH = "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf"
CUSTODY_BNB = "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88"
CUSTODY_SOL = "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm"
)

func applyFee(a float64) (float64, float64) {
fee := a * 0.003
return a - fee, fee
}

func main() {
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
"address": "EKI_31e3da3ae96a0af3",
"chainId": 701,
"name":    "EKIChain",
"status":  "Online",
"rpc":     "https://eki-chain.onrender.com",
"symbol":  "EKI",
})
})

http.HandleFunc("/bridge/deposit", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
toUser, fee := applyFee(1.0)
json.NewEncoder(w).Encode(map[string]interface{}{
"deposit":     1.0,
"toUser":      toUser,
"fee":         fee,
"feeWallet":   FEE_WALLET,
"custody_BTC": CUSTODY_BTC,
})
})

log.Println("EKI Chain 701 rodando na 8080")
log.Fatal(http.ListenAndServe(":8080", nil))
}
