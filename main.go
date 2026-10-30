package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"os"
)

func main(){
  port := os.Getenv("PORT")
  if port == "" { port = "8080" }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin","*")
    w.Header().Set("Content-Type","application/json")
    if r.Method == "POST" {
      var req map[string]interface{}
      json.NewDecoder(r.Body).Decode(&req)
      method := req["method"]
      id := req["id"]
      var result interface{} = "0x0"
      if method == "eth_chainId" { result = "0x2bd" }
      if method == "eth_blockNumber" { result = "0x1" }
      if method == "net_version" { result = "701" }
      json.NewEncoder(w).Encode(map[string]interface{}{
        "jsonrpc":"2.0","id":id,"result":result,
      })
      return
    }
    w.Write([]byte(`{"name":"EKIChain","chainId":701,"symbol":"EKI","rpc":"https://eki-chain.onrender.com","status":"Online"}`))
  })

  fmt.Printf("EKI Chain 701 RPC Live :%s\n", port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
