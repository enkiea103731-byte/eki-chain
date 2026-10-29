package main

import (
"fmt"
"log"
"net/http"
"os"
"time"
)

func main(){
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  fmt.Println("🔥 EKIChain 701 - REFEITA NO TERMUX 🔥")
  fmt.Printf("💎 ChainID 701 | Symbol EKI | Port %s\n", port)

  go func(){
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
      w.Header().Set("Content-Type","application/json")
      w.Write([]byte(`{"name":"EKIChain","chainId":701,"symbol":"EKI","status":"Online","owner":"Enki71"}`))
    })
    log.Printf("API escutando :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
  }()

  go func(){
    log.Println("P2P escutando :6001")
    http.ListenAndServe(":6001", nil)
  }()

  for { time.Sleep(1*time.Hour) }
}
