package main

import (
  "encoding/json"
  "fmt"
  "math/big"
  "net/http"
  "os"
  "strconv"
  "sync"
  "time"
)

type Block struct {
  Index int `json:"index"`
  From string `json:"from"`
  To string `json:"to"`
  Value float64 `json:"value"`
  Time string `json:"time"`
  Hash string `json:"hash,omitempty"`
}
type Stake struct {
  Wallet string `json:"wallet"`
  Amount float64 `json:"amount"`
  Until int64 `json:"until"`
}

var mu sync.Mutex
var chain []Block
var balances = map[string]float64{}
var stakes = []Stake{}
var lastMine = map[string]int64{}
var totalSupply float64
var KEY = "EKI12026SEGREDO"
const MAX_SUPPLY = 21000000.0

func save(){
  b,_ := json.Marshal(map[string]any{"balances":balances,"chain":chain,"stakes":stakes,"totalSupply":totalSupply})
  os.WriteFile("data.json", b, 0644)
}
func load(){
  b, err := os.ReadFile("data.json")
  if err != nil { return }
  var d struct{
    Balances map[string]float64 `json:"balances"`
    Chain []Block `json:"chain"`
    Stakes []Stake `json:"stakes"`
    TotalSupply float64 `json:"totalSupply"`
  }
  if json.Unmarshal(b,&d)!=nil { return }
  if d.Balances!=nil { balances=d.Balances }
  if d.Chain!=nil { chain=d.Chain }
  stakes=d.Stakes
  totalSupply=d.TotalSupply
}
func write(w http.ResponseWriter, v any){
  w.Header().Set("Content-Type","application/json")
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Methods","GET,POST,OPTIONS")
  w.Header().Set("Access-Control-Allow-Headers","Content-Type")
  json.NewEncoder(w).Encode(v)
}
func appendBlock(from, to string, val float64) Block {
  b := Block{Index: len(chain), From: from, To: to, Value: val, Time: time.Now().Format(time.RFC3339), Hash: fmt.Sprintf("0x%x", time.Now().UnixNano())}
  chain = append(chain, b)
  save()
  return b
}
func rewardAtual() float64 {
  h := len(chain)/100000
  r := 10.0
  for i:=0;i<h;i++ { r=r/2 }
  return r
}

func main(){
  load()
  if len(chain)==0 && len(balances)==0 {
    balances["enki71"]=1000010
    balances["ekifee"]=0
    totalSupply=1000010
    appendBlock("genesis","enki71",1000000)
    appendBlock("mineracao","enki71",10)
  }
  if _,ok:=balances["ekifee"]; !ok { balances["ekifee"]=0 }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    if r.Method=="OPTIONS" { write(w, map[string]any{"ok":true}); return }
    if r.Method=="POST" {
      var req map[string]any
      json.NewDecoder(r.Body).Decode(&req)
      method,_:=req["method"].(string)
      var result any="0x0"
      mu.Lock()
      switch method {
      case "eth_chainId": result="0x111d5"
      case "net_version": result="70101"
      case "web3_clientVersion": result="EKI Chain/1.0.0"
      case "net_listening": result=true
      case "eth_blockNumber": result=fmt.Sprintf("0x%x",len(chain))
      case "eth_gasPrice": result="0x0"
      case "eth_maxPriorityFeePerGas": result="0x0"
      case "eth_feeHistory": result=map[string]any{"baseFeePerGas":[]string{"0x0"},"gasUsedRatio":[]float64{0},"reward":[]any{[]string{"0x0"}},"oldestBlock":"0x0"}
      case "eth_getBalance":
        if p,ok:=req["params"].([]any); ok && len(p)>0 {
          addr,_:=p[0].(string)
          bal:=balances[addr]
          wei:=big.NewInt(int64(bal*1000000)); wei.Mul(wei,big.NewInt(1000000000000))
          result="0x"+wei.Text(16)
        }
      case "eth_getTransactionCount": result="0x0"
      case "eth_getCode": result="0x"
      case "eth_call": result="0x"
      case "eth_estimateGas": result="0x5208"
      case "eth_getBlockByNumber": result=map[string]any{"number":fmt.Sprintf("0x%x",len(chain)),"hash":fmt.Sprintf("0x%x",time.Now().UnixNano()),"timestamp":fmt.Sprintf("0x%x",time.Now().Unix()),"transactions":[]any{}}
      }
      mu.Unlock()
      write(w, map[string]any{"jsonrpc":"2.0","id":req["id"],"result":result})
      return
    }
    mu.Lock(); defer mu.Unlock()
    write(w, map[string]any{"ok":true,"chain":"EKI Chain","chainId":70101,"symbol":"EKI","blocks":len(chain),"totalSupply":totalSupply,"maxSupply":MAX_SUPPLY,"reward":rewardAtual(),"feeWallet":"ekifee","feesCollected":balances["ekifee"]})
  })

  http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    wallet:=r.URL.Query().Get("wallet"); if wallet=="" { wallet=r.URL.Query().Get("address") }
    write(w, map[string]any{"ok":true,"wallet":wallet,"balance":balances[wallet],"staked":func() float64{ s:=0.0; for _,st:=range stakes{ if st.Wallet==wallet{ s+=st.Amount } }; return s}()})
  })
  http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    write(w, map[string]any{"ok":true,"blocks":chain})
  })
  http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    write(w, map[string]any{"ok":true,"chain":"EKI Chain","chainId":70101,"blocks":len(chain),"totalSupply":totalSupply,"maxSupply":MAX_SUPPLY,"reward":rewardAtual(),"halving":len(chain)/100000,"holders":len(balances),"feesCollected":balances["ekifee"],"price":"0.10"})
  })
  http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    wallet:=r.URL.Query().Get("wallet"); if wallet=="" { wallet=r.URL.Query().Get("address") }
    if wallet=="" { write(w, map[string]any{"ok":false,"error":"falta wallet"}); return }
    now:=time.Now().Unix()
    if now-lastMine[wallet] < 60 { write(w, map[string]any{"ok":false,"error":"espere 60s para minerar de novo","nextIn":60-(now-lastMine[wallet])}); return }
    reward:=rewardAtual()
    if totalSupply+reward > MAX_SUPPLY { write(w, map[string]any{"ok":false,"error":"supply maximo atingido"}); return }
    lastMine[wallet]=now
    balances[wallet]+=reward; totalSupply+=reward
    b:=appendBlock("mineracao",wallet,reward)
    write(w, map[string]any{"ok":true,"block":b.Index,"miner":wallet,"reward":reward,"balance":balances[wallet],"totalSupply":totalSupply})
  })
  http.HandleFunc("/faucet", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    wallet:=r.URL.Query().Get("wallet"); if wallet=="" { wallet=r.URL.Query().Get("address") }
    if wallet=="" { write(w, map[string]any{"ok":false,"error":"falta wallet"}); return }
    now:=time.Now().Unix()
    if now-lastMine[wallet] < 86400 { write(w, map[string]any{"ok":false,"error":"faucet 1x por dia"}); return }
    lastMine[wallet]=now
    reward:=rewardAtual()
    balances[wallet]+=reward; totalSupply+=reward
    b:=appendBlock("mineracao",wallet,reward)
    write(w, map[string]any{"ok":true,"block":b.Index,"reward":reward,"balance":balances[wallet]})
  })
  http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    if r.URL.Query().Get("key") != KEY { write(w, map[string]any{"ok":false,"error":"sem autorizacao"}); return }
    from:=r.URL.Query().Get("from"); to:=r.URL.Query().Get("to")
    val,_:=strconv.ParseFloat(r.URL.Query().Get("value"),64)
    if to=="" || val<=0 { write(w, map[string]any{"ok":false,"error":"dados invalidos"}); return }
    fee:=val*0.01; total:=val+fee
    if balances[from] < total { write(w, map[string]any{"ok":false,"error":"saldo insuficiente, precisa de taxa 1%"}); return }
    balances[from]-=total; balances[to]+=val; balances["ekifee"]+=fee
    b:=appendBlock(from,to,val)
    write(w, map[string]any{"ok":true,"block":b.Index,"to":to,"value":val,"fee":fee,"balance":balances[to],"feesCollected":balances["ekifee"]})
  })
  http.HandleFunc("/stake", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    if r.URL.Query().Get("key") != KEY { write(w, map[string]any{"ok":false,"error":"sem autorizacao"}); return }
    wallet:=r.URL.Query().Get("wallet"); val,_:=strconv.ParseFloat(r.URL.Query().Get("value"),64)
    if balances[wallet] < val || val<=0 { write(w, map[string]any{"ok":false,"error":"saldo insuficiente"}); return }
    balances[wallet]-=val
    stakes=append(stakes, Stake{Wallet:wallet,Amount:val,Until:time.Now().Add(30*24*time.Hour).Unix()})
    save()
    write(w, map[string]any{"ok":true,"staked":val,"until":time.Now().Add(30*24*time.Hour).Format(time.RFC3339),"balance":balances[wallet]})
  })
  http.HandleFunc("/unstake", func(w http.ResponseWriter, r *http.Request){
    mu.Lock(); defer mu.Unlock()
    wallet:=r.URL.Query().Get("wallet")
    now:=time.Now().Unix(); returned:=0.0
    var rest []Stake
    for _,st:=range stakes {
      if st.Wallet==wallet && now>=st.Until {
        reward:=st.Amount*0.12/12 // 12% ao ano / 12 = 1% por 30 dias
        balances[wallet]+=st.Amount+reward; returned+=st.Amount+reward
      } else { rest=append(rest,st) }
    }
    stakes=rest; save()
    write(w, map[string]any{"ok":true,"returned":returned,"balance":balances[wallet]})
  })

  port:=os.Getenv("PORT"); if port=="" { port="10000" }
  fmt.Println("EKI CHAIN PRO rodando na porta", port, "reward", rewardAtual(), "supply", totalSupply)
  http.ListenAndServe(":"+port, nil)
}
