import json, time, os
from http.server import BaseHTTPRequestHandler, HTTPServer
BLOCK = 70101
START = int(time.time())
BAL_FILE='balances.json'
def load():
    try:
        with open(BAL_FILE) as f: return json.load(f)
    except: return {"balances":{"0x5d3f7f34ea94fe40ae783a258e246e293cbac206":1000000},"stakes":{},"total_supply":1000000,"claimed":[]}
DATA=load()
def save():
    with open(BAL_FILE,'w') as f: json.dump(DATA,f)
def wei(x): return hex(int(x*10**18))
def read_index():
    try:
        with open('index.html','rb') as f: return f.read()
    except: return b'<h1>EKI CHAIN MAINNET ONLINE - 70101</h1>'
class H(BaseHTTPRequestHandler):
    def log_message(self,*a): pass
    def send_json(self, obj, code=200):
        d=json.dumps(obj).encode()
        self.send_response(code); self.send_header('Content-Type','application/json'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Content-Length',str(len(d))); self.end_headers(); self.wfile.write(d)
    def send_html(self):
        data=read_index()
        self.send_response(200); self.send_header('Content-Type','text/html; charset=utf-8'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Cache-Control','no-cache'); self.send_header('Content-Length',str(len(data))); self.end_headers(); self.wfile.write(data)
    def do_GET(self):
        p=self.path.split('?')[0]
        if p=='/stats':
            return self.send_json({'chain':'EKI Chain MAINNET','chainId':70101,'block':BLOCK+int((time.time()-START)//5),'total_supply':DATA['total_supply'],'holders':len(DATA['balances']),'circulating':sum(DATA['balances'].values()),'staked_total':sum(DATA['stakes'].values())})
        if p=='/health' or p=='/ping':
            self.send_response(200); self.send_header('Content-Length','2'); self.end_headers(); self.wfile.write(b'ok'); return
        if p.startswith('/balance/'):
            a=p.split('/balance/')[-1].lower()
            return self.send_json({'address':a,'balance':DATA['balances'].get(a,0),'staked':DATA['stakes'].get(a,0)})
        return self.send_html()
    def do_POST(self):
        l=int(self.headers.get('Content-Length',0)); body=self.rfile.read(l) if l else b'{}'
        try: req=json.loads(body or b'{}')
        except: req={}
        path=self.path.split('?')[0]
        if path=='/claim':
            a=req.get('address','').lower()
            if not a.startswith('0x') or len(a)!=42: return self.send_json({'error':'endereco invalido'},400)
            if a in DATA['claimed']: return self.send_json({'error':'voce ja pegou, 1x por carteira'},400)
            DATA['balances'][a]=DATA['balances'].get(a,0)+10; DATA['claimed'].append(a); save()
            return self.send_json({'ok':True,'balance':DATA['balances'][a],'msg':'Airdrop MAINNET: 10 $EKI'})
        if path=='/stake':
            a=req.get('address','').lower(); amt=float(req.get('amount',0))
            if DATA['balances'].get(a,0)<amt: return self.send_json({'error':'saldo insuficiente'},400)
            DATA['balances'][a]-=amt; DATA['stakes'][a]=DATA['stakes'].get(a,0)+amt; save()
            return self.send_json({'ok':True,'staked':DATA['stakes'][a],'apy':'12% ao ano'})
        if path=='/transfer':
            f=req.get('from','').lower(); t=req.get('to','').lower(); amt=float(req.get('amount',0))
            if DATA['balances'].get(f,0)<amt: return self.send_json({'error':'saldo insuficiente'},400)
            DATA['balances'][f]-=amt; DATA['balances'][t]=DATA['balances'].get(t,0)+amt; save()
            return self.send_json({'ok':True,'from_balance':DATA['balances'][f],'to_balance':DATA['balances'][t]})
        method=req.get('method',''); rid=req.get('id',1); params=req.get('params',[])
        cur=BLOCK+int((time.time()-START)//5); res=None
        if method=='eth_chainId': res='0x111d5'
        elif method=='net_version': res='70101'
        elif method=='web3_clientVersion': res='EKI Chain MAINNET/1.0.0'
        elif method=='eth_blockNumber': res=hex(cur)
        elif method=='eth_gasPrice': res='0x0'
        elif method=='eth_getBalance': res=wei(DATA['balances'].get((params[0] if params else '').lower(),0))
        elif method=='net_listening': res=True
        elif method=='eth_getBlockByNumber': res={'number':hex(cur),'hash':'0x'+'0'*64,'timestamp':hex(int(time.time())),'transactions':[]}
        else: res='0x0'
        return self.send_json({'jsonrpc':'2.0','id':rid,'result':res})
    def do_OPTIONS(self):
        self.send_response(200); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Access-Control-Allow-Methods','GET,POST,OPTIONS'); self.end_headers()
PORT=int(os.environ.get('PORT',10000))
print(f"EKI MAINNET na porta {PORT}", flush=True)
HTTPServer(('0.0.0.0',PORT), H).serve_forever()
