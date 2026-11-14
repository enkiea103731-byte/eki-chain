package main
import ("crypto/sha256";"encoding/hex";"encoding/json";"fmt";"net/http";"os";"sync";"time")
type Block struct{Index int;Timestamp string;Data string;PrevHash string;Hash string}
type Chain struct{Blocks []Block;Balances map[string]float64}
var mu sync.Mutex
var bc *Chain
func hash(b Block) string {r:=fmt.Sprintf("%d%s%s%s",b.Index,b.Timestamp,b.Data,b.PrevHash);h:=sha256.New();h.Write([]byte(r));return hex.EncodeToString(h.Sum(nil))}
func load() *Chain {c:=&Chain{Blocks:[]Block{},Balances:map[string]float64{"enki71":1000000}};d,_:=os.ReadFile("eki.json");if len(d)>0{json.Unmarshal(d,c)};if len(c.Blocks)==0{g:=Block{0,time.Now().String(),"GENESIS $EKI","",""};g.Hash=hash(g);c.Blocks=append(c.Blocks,g)};return c}
func save(){d,_:=json.Marshal(bc);os.WriteFile("eki.json",d,0644)}
func miner(){for{time.Sleep(2*time.Second);mu.Lock();prev:=bc.Blocks[len(bc.Blocks)-1];nb:=Block{prev.Index+1,time.Now().String(),fmt.Sprintf("Bloco %d",prev.Index+1),prev.Hash,""};nb.Hash=hash(nb);bc.Blocks=append(bc.Blocks,nb);bc.Balances["enki71"]+=10;save();fmt.Printf("Block %d Saldo %.2f\n",nb.Index,bc.Balances["enki71"]);mu.Unlock()}}
func main(){
bc=load();go miner()
http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){mu.Lock();defer mu.Unlock();w.Header().Set("Content-Type","application/json");json.NewEncoder(w).Encode(map[string]interface{}{"chain":"EKI","blocks":len(bc.Blocks),"last_hash":bc.Blocks[len(bc.Blocks)-1].Hash,"balances":bc.Balances})})
http.HandleFunc("/balance",func(w http.ResponseWriter,r *http.Request){mu.Lock();defer mu.Unlock();name:=r.URL.Query().Get("name");if name==""{name="enki71"};json.NewEncoder(w).Encode(map[string]interface{}{"name":name,"balance":bc.Balances[name]})})
http.HandleFunc("/send",func(w http.ResponseWriter,r *http.Request){mu.Lock();defer mu.Unlock();from:=r.URL.Query().Get("from");to:=r.URL.Query().Get("to");var v float64;fmt.Sscanf(r.URL.Query().Get("value"),"%f",&v);if bc.Balances[from]<v{json.NewEncoder(w).Encode(map[string]string{"error":"saldo insuficiente"});return};bc.Balances[from]-=v;bc.Balances[to]+=v;prev:=bc.Blocks[len(bc.Blocks)-1];nb:=Block{prev.Index+1,time.Now().String(),fmt.Sprintf("TX %s -> %s %.2f",from,to,v),prev.Hash,""};nb.Hash=hash(nb);bc.Blocks=append(bc.Blocks,nb);save();json.NewEncoder(w).Encode(map[string]interface{}{"ok":true,"block":nb.Index})})
port:=os.Getenv("PORT");if port==""{port="8080"};fmt.Println("EKI na porta",port);http.ListenAndServe(":"+port,nil)
}
