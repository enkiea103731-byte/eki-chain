package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "os"
  "strconv"
  "sync"
  "time"
)

type Block struct {
  Index int     `json:"index"`
  From  string  `json:"from"`
  To    string  `json:"to"`
  Value float64 `json:"value"`
  Time  string  `json:"time"`
}

var mu sync.Mutex
var balances = map[string]float64{"enki71": 1000000}
var chain []Block
const KEY = "ENKI2026SEGREDO"

func write(w http.ResponseWriter, v any) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  json.NewEncoder(w).Encode(v)
}

func main() {
  chain = append(chain, Block{Index: 0, From: "genesis", To: "enki71", Value: 1000000, Time: time.Now().Format(time.RFC3339)})

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    write(w, map[string]any{"ok": true, "coin": "EKI", "blocks": len(chain)})
  })

  http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); defer mu.Unlock()
    wallet := r.URL.Query().Get("wallet")
    if wallet == "" { wallet = r.URL.Query().Get("address") }
    write(w, map[string]any{"ok": true, "wallet": wallet, "balance": balances[wallet]})
  })

  http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); defer mu.Unlock()
    write(w, map[string]any{"ok": true, "blocks": chain})
  })

  http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); defer mu.Unlock()
    wallet := r.URL.Query().Get("wallet")
    if wallet == "" { wallet = r.URL.Query().Get("address") }
    if wallet == "" { wallet = r.URL.Query().Get("miner") }
    if wallet == "" {
      write(w, map[string]any{"ok": false, "error": "falta wallet"})
      return
    }
    reward := 10.0
    balances[wallet] += reward
    b := Block{Index: len(chain), From: "mineracao", To: wallet, Value: reward, Time: time.Now().Format(time.RFC3339)}
    chain = append(chain, b)
    write(w, map[string]any{"ok": true, "block": b.Index, "miner": wallet, "reward": reward, "balance": balances[wallet]})
  })

  http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); defer mu.Unlock()
    if r.URL.Query().Get("key") != KEY {
      write(w, map[string]any{"ok": false, "error": "sem autorizacao"})
      return
    }
    from := r.URL.Query().Get("from")
    to := r.URL.Query().Get("to")
    val, _ := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
    if to == "" || val <= 0 {
      write(w, map[string]any{"ok": false, "error": "dados invalidos"})
      return
    }
    if balances[from] < val {
      write(w, map[string]any{"ok": false, "error": "saldo insuficiente"})
      return
    }
    balances[from] -= val
    balances[to] += val
    b := Block{Index: len(chain), From: from, To: to, Value: val, Time: time.Now().Format(time.RFC3339)}
    chain = append(chain, b)
    write(w, map[string]any{"ok": true, "block": b.Index, "to": to, "value": val, "balance": balances[to]})
  })

  port := os.Getenv("PORT")
  if port == "" { port = "10000" }
  fmt.Println("EKI rodando na porta", port)
  http.ListenAndServe(":"+port, nil)
}
