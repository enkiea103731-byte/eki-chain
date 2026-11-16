package main

import (
	"fmt"
	"net/http"
)

func hasRPC(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func init() {
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			return
		}
		b := make([]byte, 8192)
		n, _ := r.Body.Read(b)
		s := string(b[:n])
		if hasRPC(s, "net_version") {
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"70101"}`)
			return
		}
		if hasRPC(s, "eth_blockNumber") {
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"0x1"}`)
			return
		}
		if hasRPC(s, "eth_gasPrice") {
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"0x0"}`)
			return
		}
		if hasRPC(s, "eth_getBalance") {
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"0x0"}`)
			return
		}
		if hasRPC(s, "eth_getBlockByNumber") {
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":{"number":"0x1","hash":"0x000279719ef7","parentHash":"0x000000000000","timestamp":"0x1","gasLimit":"0x1c9c380","gasUsed":"0x0","transactions":[]}}`)
			return
		}
		fmt.Fprint(w, `{"jsonrpc":"2.0","id":1,"result":"0x111d5"}`)
	})
}
