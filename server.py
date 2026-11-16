import json, time, os
from http.server import BaseHTTPRequestHandler, HTTPServer

BLOCK = 70101
START = int(time.time())

class H(BaseHTTPRequestHandler):
    def log_message(self, *a): pass
    def do_GET(self):
        if self.path in ['/', '/index.html']:
            try:
                with open('index.html','rb') as f: data=f.read()
            except: data=b'<h1>EKI CHAIN ONLINE</h1>'
            self.send_response(200); self.send_header('Content-Type','text/html'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Content-Length',str(len(data))); self.end_headers(); self.wfile.write(data)
        elif self.path in ['/health','/ping']:
            d=b'ok'; self.send_response(200); self.send_header('Content-Length','2'); self.end_headers(); self.wfile.write(d)
        else:
            self.send_response(404); self.end_headers()
    def do_POST(self):
        length=int(self.headers.get('Content-Length',0))
        body=self.rfile.read(length) if length else b'{}'
        try: req=json.loads(body or b'{}')
        except: req={}
        method=req.get('method',''); rid=req.get('id',1); params=req.get('params',[])
        # bloco sobe 1 a cada 5 segundos
        cur_block = BLOCK + int((time.time()-START)//5)
        result=None
        if method=='eth_chainId': result='0x111d5'
        elif method=='net_version': result='70101'
        elif method=='web3_clientVersion': result='EKI Chain/1.0.0'
        elif method=='eth_blockNumber': result=hex(cur_block)
        elif method=='eth_gasPrice': result='0x0'
        elif method=='eth_getBalance': result='0x0'
        elif method=='eth_accounts': result=[]
        elif method=='net_listening': result=True
        elif method=='eth_getBlockByNumber':
            n=params[0] if params else 'latest'
            result={'number':hex(cur_block),'hash':'0x'+'0'*64,'parentHash':'0x'+'0'*64,'timestamp':hex(int(time.time())),'gasLimit':'0x7a1200','gasUsed':'0x0','miner':'0x5d3f7f34ea94fe40ae783a258e246e293cbac206','transactions':[],'uncles':[]}
        elif method=='eth_getBlockByHash': result={'number':hex(cur_block),'hash':params[0] if params else '0x0','transactions':[]}
        else: result='0x0'
        resp=json.dumps({'jsonrpc':'2.0','id':rid,'result':result}).encode()
        self.send_response(200); self.send_header('Content-Type','application/json'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Content-Length',str(len(resp))); self.end_headers(); self.wfile.write(resp)
    def do_OPTIONS(self):
        self.send_response(200); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Access-Control-Allow-Methods','GET,POST,OPTIONS'); self.end_headers()

PORT=int(os.environ.get('PORT',10000))
print(f"EKI RPC na porta {PORT}", flush=True)
HTTPServer(('0.0.0.0',PORT), H).serve_forever()
