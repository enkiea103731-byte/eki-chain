package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const SENHA_CONFIRMAR = "EKI-MARTA-20-REAL"
const PRECO = 0.10

type Pedido struct {
	ID     int
	Nome   string
	Valor  float64
	EKI    float64
	Status string
}

type Chain struct {
	Blocos    int
	SaldoEnki float64
	Saldos    map[string]float64
	Pedidos   []Pedido
	NextID    int
}

var mu sync.Mutex
var c Chain

func load() {
	data, err := os.ReadFile("chain.json")
	if err == nil {
		json.Unmarshal(data, &c)
	}
	if c.Saldos == nil {
		c.Saldos = map[string]float64{}
	}
	if c.NextID == 0 {
		c.NextID = 1
	}
	// LIMPA FRAUDE: mantem so a Marta que pagou de verdade
	var limpos []Pedido
	for _, p := range c.Pedidos {
		if p.Nome == "Marta" && p.ID == 1 {
			limpos = append(limpos, p)
		}
	}
	c.Pedidos = limpos
	// Garante saldo da Marta real
	c.Saldos["Marta"] = 200
}

func save() {
	data, _ := json.Marshal(c)
	os.WriteFile("chain.json", data, 0644)
}

func minerar() {
	for {
		time.Sleep(2 * time.Second)
		mu.Lock()
		c.Blocos++
		c.SaldoEnki += 10
		c.Saldos["enki71"] = c.SaldoEnki
		save()
		mu.Unlock()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	load()
	go minerar()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>EKI - REAL</h1>")
		fmt.Fprintf(w, "Blocos: %d | Minerando 24/7 | Dificuldade 3 | Preco: R$ %.2f por EKI<br><br>", c.Blocos, PRECO)
		fmt.Fprintf(w, "PIX: <b>32468719816</b><br><br>")
		fmt.Fprintf(w, "Saldo enki71: %.2f EKI<hr>", c.SaldoEnki)
		fmt.Fprintf(w, `<h2>COMPRAR</h2><form action="/buy" method="GET">Nome: <input name="nome"><br>Valor R$: <input name="valor" type="number" step="0.01"><br><input type="submit" value="GERAR PEDIDO PIX"></form><hr>`)
		fmt.Fprintf(w, "<h2>PEDIDOS</h2>")
		for _, p := range c.Pedidos {
			fmt.Fprintf(w, "#%d %s - R$ %.2f = %.2f EKI - %s<br>", p.ID, p.Nome, p.Valor, p.EKI, p.Status)
			// SEM BOTAO CONFIRMAR - NINGUEM MAIS CONSEGUE CLICAR
		}
		fmt.Fprintf(w, `<br><br><a href="/balance?nome=Marta">Ver saldo Marta</a>`)
	})

	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		nome := r.URL.Query().Get("nome")
		valor, _ := strconv.ParseFloat(r.URL.Query().Get("valor"), 64)
		if nome == "" || valor <= 0 {
			http.Redirect(w, r, "/", 302)
			return
		}
		eki := valor / PRECO
		p := Pedido{ID: c.NextID, Nome: nome, Valor: valor, EKI: eki, Status: "AGUARDANDO PIX"}
		c.NextID++
		c.Pedidos = append(c.Pedidos, p)
		save()
		fmt.Fprintf(w, "Pedido #%d criado! %s pague R$ %.2f no PIX 32468719816 para receber %.2f EKI. Aguarde confirmacao do dono.", p.ID, nome, valor, eki)
	})

	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		// TRAVA DE SEGURANCA - SO COM SENHA
		key := r.URL.Query().Get("key")
		if key != SENHA_CONFIRMAR {
			http.Error(w, "Sem permissao - fraude bloqueada", 403)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		for i, p := range c.Pedidos {
			if p.ID == id {
				if p.Status == "PAGO" {
					fmt.Fprintf(w, "Ja estava pago")
					return
				}
				c.Pedidos[i].Status = "PAGO"
				c.Saldos[p.Nome] += p.EKI
				if c.SaldoEnki >= p.EKI {
					c.SaldoEnki -= p.EKI
				}
				save()
				fmt.Fprintf(w, "Confirmado! %s recebeu %.2f EKI", p.Nome, p.EKI)
				return
			}
		}
		fmt.Fprintf(w, "Pedido nao encontrado")
	})

	http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		nome := r.URL.Query().Get("nome")
		fmt.Fprintf(w, "Saldo de %s: %.2f EKI", nome, c.Saldos[nome])
	})

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"0x111d5"}`)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	fmt.Println("Site no ar na porta", port, "- Minerando 24/7 - Anti-fraude ON")
	http.ListenAndServe(":"+port, nil)
}
