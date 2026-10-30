package main
import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "io"
 "log"
 "net/http"
 "strings"
 "time"
)
const FIXED_ADDR = "EKI_31e3da3ae96a0af3"
const FIXED_PRIV = "31e3da3ae96a0af3b1c2d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6"
type Block struct{Index int; Time,Data,Prev,Hash string}
var Chain []Block
func main(){
 Chain=append(Chain, Block{0,time.Now().String(),"GENESIS EKI 21M 701 FIXED 31e3","",""})
 h:=sha256.Sum256([]byte(fmt.Sprint(Chain[0]))); Chain[0].Hash=hex.EncodeToString(h[:])
 fmt.Println("Carteira FIXA:",FIXED_ADDR)
 http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
  if r.Method=="POST"{body,_:=io.ReadAll(r.Body); s:=string(body); w.Header().Set("Content-Type","application/json"); w.Header().Set("Access-Control-Allow-Origin","*")
   if strings.Contains(s,"eth_chainId"){w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x2bd"}`)); return}
   if strings.Contains(s,"net_version"){w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"701"}`)); return}
   if strings.Contains(s,"eth_blockNumber"){w.Write([]byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":"0x%x"}`,len(Chain)))); return}
   w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x0"}`)); return
  }
  w.Header().Set("Access-Control-Allow-Origin","*"); json.NewEncoder(w).Encode(map[string]interface{}{"chainId":701,"name":"EKIChain","symbol":"EKI","rpc":"https://eki-chain.onrender.com","status":"Online","address":FIXED_ADDR})
 })
 log.Fatal(http.ListenAndServe(":8080",nil))
}
