package main
import (
 "crypto/rand"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "io"
 "log"
 "net/http"
 "os"
 "strings"
 "time"
)
type Wallet struct{PrivateKey,PublicKey,Address string}
type Block struct{Index int; Time,Data,Prev,Hash string}
var Chain []Block
var W Wallet
func main(){
 if d,err:=os.ReadFile("wallet.json"); err==nil { json.Unmarshal(d,&W) } else {
  b:=make([]byte,32); rand.Read(b)
  pk:=hex.EncodeToString(b)
  W=Wallet{pk,pk[:32],"EKI_"+pk[:16]}
  j,_:=json.Marshal(W); os.WriteFile("wallet.json",j,0644)
 }
 Chain=append(Chain, Block{0,time.Now().String(),"GENESIS EKI 21M 701","",""})
 h:=sha256.Sum256([]byte(fmt.Sprint(Chain[0]))); Chain[0].Hash=hex.EncodeToString(h[:])
 fmt.Println("Carteira:",W.Address)
 http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
  if r.Method=="POST"{body,_:=io.ReadAll(r.Body); s:=string(body); w.Header().Set("Content-Type","application/json"); w.Header().Set("Access-Control-Allow-Origin","*")
   if strings.Contains(s,"eth_chainId"){w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x2bd"}`)); return}
   if strings.Contains(s,"net_version"){w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"701"}`)); return}
   if strings.Contains(s,"eth_blockNumber"){w.Write([]byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":"0x%x"}`,len(Chain)))); return}
   w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x0"}`)); return
  }
  w.Header().Set("Access-Control-Allow-Origin","*"); json.NewEncoder(w).Encode(map[string]interface{}{"name":"EKIChain","chainId":701,"symbol":"EKI","rpc":"https://eki-chain.onrender.com","status":"Online","address":W.Address})
 })
 http.HandleFunc("/chain",func(w http.ResponseWriter, r *http.Request){w.Header().Set("Access-Control-Allow-Origin","*"); json.NewEncoder(w).Encode(Chain)})
 fmt.Println("API :8080 - Chain 701 Online")
 log.Fatal(http.ListenAndServe(":8080",nil))
}
