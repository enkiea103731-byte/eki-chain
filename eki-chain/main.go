package main

import (
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "os"
 "sync"
 "time"
)

type State struct {
 Blocks int `json:"blocks"`
 Balances map[string]float64 `json:"balances"`
}

var mu sync.Mutex
var st State

func load() {
 b, _ := os.ReadFile("eki.json")
 if len(b) > 0 { json.Unmarshal(b, &st) }
 if st.Balances == nil { st.Balances = map[string]float64{"enki71": 1000000} }
}
func save() {
 b, _ := json.Marshal(st)
 os.WriteFile("eki.json", b, 0644)
}

func main() {
 load()
 go func() {
  for {
   time.Sleep(10*time.Second)
   mu.Lock()
   if st.Balances["enki71"] < 21000000 {
    st.Blocks++
    st.Balances["enki71"] += 10
    save()
   }
   mu.Unlock()
  }
 }()

 // CORS pra qualquer carteira
 cors := func(w http.ResponseWriter) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  w.Header().Set("Access-Control-Allow-Methods", "*")
 }

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cors(w); w.Header().Set("Content-Type","application/json")
  mu.Lock(); defer mu.Unlock()
  json.NewEncoder(w).Encode(st)
 })

 http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
  cors(w); w.Header().Set("Content-Type","application/json")
  if r.Method == "OPTIONS" { return }
  var req struct{ Method string `json:"method"`; ID any `json:"id"`; Params []any `json:"params"` }
  json.NewDecoder(r.Body).Decode(&req)
  mu.Lock(); defer mu.Unlock()
  var result any = "0x0"
  switch req.Method {
  case "eth_chainId", "net_version":
   result = "0x111d5"
   if req.Method=="net_version" { result="70101" }
  case "eth_blockNumber":
   result = fmt.Sprintf("0x%x", st.Blocks)
  case "eth_getBalance":
   result = "0xde0b6b3a7640000" // 1 EKI exemplo
  case "eth_gasPrice":
   result = "0x0"
  case "eth_getTransactionCount":
   result = "0x0"
  }
  json.NewEncoder(w).Encode(map[string]any{"jsonrpc":"2.0","id":req.ID,"result":result})
 })

 port := os.Getenv("PORT")
 if port=="" { port="8080" }
 log.Println("EKI 70101 vivo na porta", port)
 log.Fatal(http.ListenAndServe(":"+port, nil))
}
