package main

import (
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "net/http"
  "os"
  "strings"
  "sync"
  "time"
)

type Block struct {
  Index int `json:"index"`
  Timestamp string `json:"time"`
  Data string `json:"data"`
  PrevHash string `json:"prev"`
  Hash string `json:"hash"`
  Nonce int `json:"nonce"`
}
type Pedido struct {
  ID int `json:"id"`
  Nome string `json:"nome"`
  Valor float64 `json:"valor"`
  EKI float64 `json:"eki"`
  Status string `json:"status"`
}
type Chain struct {
  Blocks []Block `json:"blocks"`
  Balances map[string]float64 `json:"balances"`
  Pedidos []Pedido `json:"pedidos"`
}

const DIFFICULTY = 3
const REWARD = 10.0
const PRECO = 0.10
const FILE = "chain.json"
const PIX_CHAVE = "32468719816"
const ADMIN_KEY = "enki71admin"

var mu sync.Mutex

func calcHash(b Block) string {
  r := fmt.Sprintf("%d%s%s%s%d", b.Index, b.Timestamp, b.Data, b.PrevHash, b.Nonce)
  h := sha256.Sum256([]byte(r))
  return hex.EncodeToString(h[:])
}
func load() Chain {
  data, err := os.ReadFile(FILE)
  if err != nil {
    g := Block{0, time.Now().String(), "GENESIS EKI", "", "", 0}
    g.Hash = calcHash(g)
    return Chain{Blocks: []Block{g}, Balances: map[string]float64{}, Pedidos: []Pedido{}}
  }
  var c Chain
  json.Unmarshal(data, &c)
  if c.Balances == nil { c.Balances = map[string]float64{} }
  return c
}
func save(c Chain) {
  data, _ := json.Marshal(c)
  os.WriteFile(FILE, data, 0644)
}
func mineLoop() {
  for {
    mu.Lock()
    c := load()
    prev := c.Blocks[len(c.Blocks)-1]
    nb := Block{prev.Index+1, time.Now().String(), "Minerado por enki71", prev.Hash, "", 0}
    prefix := strings.Repeat("0", DIFFICULTY)
    tentativas := 0
    for {
      nb.Nonce++
      tentativas++
      nb.Hash = calcHash(nb)
      if strings.HasPrefix(nb.Hash, prefix) { break }
    }
    c.Blocks = append(c.Blocks, nb)
    c.Balances["enki71"] += REWARD
    save(c)
    fmt.Printf("Block %d MINERADO! Hash %s Nonce %d Tentativas %d Saldo enki71: %.2f\n", nb.Index, nb.Hash[:12], nb.Nonce, tentativas, c.Balances["enki71"])
    mu.Unlock()
    time.Sleep(1 * time.Second)
  }
}
func main() {
  go mineLoop()
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); c := load(); mu.Unlock()
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, `<h1>EKI - REAL</h1><p>Blocos: %d | Minerando 24/7 | Dificuldade %d | Preco: R$ %.2f por EKI</p><p>PIX: <b>%s</b></p><p>Saldo enki71: %.2f EKI</p><hr><h3>COMPRAR</h3><form action="/comprar" method="POST">Nome: <input name="nome"><br>Valor R$: <input name="valor" type="number" step="0.01"><br><button>GERAR PEDIDO PIX</button></form><hr><h3>PEDIDOS</h3>`, len(c.Blocks), DIFFICULTY, PRECO, PIX_CHAVE, c.Balances["enki71"])
    for _, p := range c.Pedidos {
      fmt.Fprintf(w, `#%d %s - R$ %.2f = %.2f EKI - %s <a href="/confirmar?id=%d&key=%s">CONFIRMAR (so voce)</a><br>`, p.ID, p.Nome, p.Valor, p.EKI, p.Status, p.ID, ADMIN_KEY)
    }
  })
  http.HandleFunc("/comprar", func(w http.ResponseWriter, r *http.Request) {
    nome := r.FormValue("nome")
    var valor float64
    fmt.Sscanf(r.FormValue("valor"), "%f", &valor)
    if nome == "" || valor <= 0 { http.Redirect(w, r, "/", 302); return }
    mu.Lock()
    c := load()
    id := len(c.Pedidos)+1
    eki := valor / PRECO
    c.Pedidos = append(c.Pedidos, Pedido{id, nome, valor, eki, "AGUARDANDO PIX"})
    save(c)
    mu.Unlock()
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, `<h1>PEDIDO #%d CRIADO</h1><p>%s, pague <b>R$ %.2f</b> no PIX <b>%s</b></p><p>Voce vai receber <b>%.2f EKI</b></p><p>Depois de pagar, me manda o comprovante no WhatsApp. Eu confirmo aqui.</p><a href="/">VOLTAR</a>`, id, nome, valor, PIX_CHAVE, eki)
  })
  http.HandleFunc("/confirmar", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("key") != ADMIN_KEY { fmt.Fprint(w, "Sem permissao"); return }
    var id int
    fmt.Sscanf(r.URL.Query().Get("id"), "%d", &id)
    mu.Lock()
    c := load()
    for i := range c.Pedidos {
      if c.Pedidos[i].ID == id && c.Pedidos[i].Status == "AGUARDANDO PIX" {
        c.Pedidos[i].Status = "PAGO"
        c.Balances[c.Pedidos[i].Nome] += c.Pedidos[i].EKI
      }
    }
    save(c)
    mu.Unlock()
    http.Redirect(w, r, "/", 302)
  })
  http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock(); c := load(); mu.Unlock()
    fmt.Fprintf(w, "Saldo de %s: %.2f EKI", r.URL.Query().Get("nome"), c.Balances[r.URL.Query().Get("nome")])
  })
  port := os.Getenv("PORT")
  if port == "" { port = "10000" }
  fmt.Println("Site no ar na porta", port, "- Minerando 24/7")
  http.ListenAndServe(":"+port, nil)
}
