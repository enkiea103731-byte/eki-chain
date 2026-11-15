package main

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"net/http"
"os"
"sync"
"time"
)

type Block struct {
Index     int     `json:"index"`
Timestamp string  `json:"timestamp"`
Data      string  `json:"data"`
PrevHash  string  `json:"prevHash"`
Hash      string  `json:"hash"`
}
type Chain struct {
Blocks   []Block            `json:"blocks"`
Balances map[string]float64 `json:"balances"`
}

var bc Chain
var mu sync.Mutex

func hash(b Block) string {
h := sha256.New()
h.Write([]byte(fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp, b.Data, b.PrevHash)))
return hex.EncodeToString(h.Sum(nil))
}
func save() { d, _ := json.Marshal(bc); os.WriteFile("eki.json", d, 0644) }
func load() {
d, err := os.ReadFile("eki.json")
if err != nil || len(d) == 0 { return }
json.Unmarshal(d, &bc)
}

func main() {
bc.Balances = map[string]float64{"enki71": 10000000}
load()
if len(bc.Blocks) == 0 {
b := Block{Index: 0, Timestamp: time.Now().String(), Data: "GENESIS EKI", PrevHash: "0"}
b.Hash = hash(b)
bc.Blocks = []Block{b}
if bc.Balances == nil { bc.Balances = map[string]float64{"enki71": 10000000} }
save()
}
if bc.Balances == nil { bc.Balances = map[string]float64{} }

// MINERAÇÃO - 10 EKI a cada 2s pra enki71
go func() {
for {
time.Sleep(2 * time.Second)
mu.Lock()
nb := Block{Index: len(bc.Blocks), Timestamp: time.Now().String(), Data: fmt.Sprintf("Bloco %d minerado", len(bc.Blocks)), PrevHash: bc.Blocks[len(bc.Blocks)-1].Hash}
nb.Hash = hash(nb)
bc.Blocks = append(bc.Blocks, nb)
bc.Balances["enki71"] += 10
save()
mu.Unlock()
fmt.Printf("Block %d - Saldo enki71: %.0f\n", nb.Index, bc.Balances["enki71"])
}
}()

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
mu.Lock(); defer mu.Unlock()
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"chain": "EKI", "blocks": len(bc.Blocks), "last_hash": bc.Blocks[len(bc.Blocks)-1].Hash, "balances": bc.Balances})
})

http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
mu.Lock(); defer mu.Unlock()
name := r.URL.Query().Get("name")
if name == "" { name = "enki71" }
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"name": name, "balance": bc.Balances[name]})
})

// TRANSFERÊNCIA - O QUE FALTAVA
http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
mu.Lock(); defer mu.Unlock()
w.Header().Set("Content-Type", "application/json")
from := r.URL.Query().Get("from")
to := r.URL.Query().Get("to")
var f float64
fmt.Sscanf(r.URL.Query().Get("value"), "%f", &f)
if from == "" { from = "enki71" }
if to == "" || f <= 0 {
json.NewEncoder(w).Encode(map[string]interface{}{"error": "use ?from=enki71&to=destino&value=10"})
return
}
if bc.Balances[from] < f {
json.NewEncoder(w).Encode(map[string]interface{}{"error": "saldo insuficiente"})
return
}
bc.Balances[from] -= f
bc.Balances[to] += f
nb := Block{Index: len(bc.Blocks), Timestamp: time.Now().String(), Data: fmt.Sprintf("TX %s->%s %.2f", from, to, f), PrevHash: bc.Blocks[len(bc.Blocks)-1].Hash}
nb.Hash = hash(nb)
bc.Blocks = append(bc.Blocks, nb)
save()
json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "block": nb.Index, "from": from, "to": to, "value": f, "novo_saldo_from": bc.Balances[from]})
})

port := os.Getenv("PORT")
if port == "" { port = "10000" }
fmt.Println("EKI rodando na porta", port)
http.ListenAndServe(":"+port, nil)
}
