package main
import (
 "fmt"
 "log"
 "net/http"
 "time"
)
func main(){
 fmt.Println("\n🔥 EKIChain 701 - REFEITA NO TERMUX 🔥")
 fmt.Println("💎 ChainID 701 | Symbol EKI | Supply 2.1B")
 fmt.Println("\n🚀 EKI Chain rodando! API :8080 | P2P :6001")
 go func(){
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
   w.Write([]byte("EKIChain 701 Online - Enki71 - API OK"))
  })
  log.Println("API escutando :8080")
  http.ListenAndServe(":8080", nil)
 }()
 go func(){
  log.Println("P2P escutando :6001")
  http.ListenAndServe(":6001", nil)
 }()
 for{ time.Sleep(1*time.Hour) }
}
