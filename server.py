import json, time, os, sqlite3, hashlib
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse

CHAIN_ID = 70101
CHAIN_HEX = '0x111D5'
GENESIS_TIME = 1763000000
START = int(time.time())
DB = 'eki_pro.db'
OWNER = '0x5d3f7f34ea94fe40ae783a258e246e293cbac206'
USER_ADDR = '0xb8ce0fe09ac0cd186a1006a4925c5af644766b45'

# --- BANCO PRO (SQLite, persiste) ---
con = sqlite3.connect(DB, check_same_thread=False)
cur = con.cursor()
cur.execute('CREATE TABLE IF NOT EXISTS balances (addr TEXT PRIMARY KEY, balance REAL, staked REAL)')
cur.execute('CREATE TABLE IF NOT EXISTS claimed (addr TEXT PRIMARY KEY)')
cur.execute('CREATE TABLE IF NOT EXISTS txs (hash TEXT PRIMARY KEY, frm TEXT, toaddr TEXT, amount REAL, ts INTEGER, block INTEGER)')
cur.execute('CREATE TABLE IF NOT EXISTS meta (k TEXT PRIMARY KEY, v TEXT)')
con.commit()

def get_bal(a):
    a=a.lower()
    r=cur.execute('SELECT balance, staked FROM balances WHERE addr=?',(a,)).fetchone()
    if r: return r
    return (0,0)
def set_bal(a,b,s):
    cur.execute('INSERT INTO balances(addr,balance,staked) VALUES(?,?,?) ON CONFLICT(addr) DO UPDATE SET balance=excluded.balance, staked=excluded.staked',(a.lower(),b,s)); con.commit()
# Genesis - igual chain de verdade
if not cur.execute('SELECT * FROM balances WHERE addr=?',(OWNER,)).fetchone():
    set_bal(OWNER, 999990, 0)
if not cur.execute('SELECT * FROM balances WHERE addr=?',(USER_ADDR,)).fetchone():
    set_bal(USER_ADDR, 10, 0)
    cur.execute('INSERT OR IGNORE INTO claimed(addr) VALUES(?)',(USER_ADDR,)); con.commit()

def block_now():
    return 70334 + int((time.time()-START)//5)
def wei(x): return hex(int(float(x)*10**18))
def to_hex(n): return hex(int(n))
def read_index():
    try:
        with open('index.html','rb') as f: return f.read()
    except: return b'<h1>EKI CHAIN MAINNET PRO - 70101</h1>'

class H(BaseHTTPRequestHandler):
    def log_message(self,*a): pass
    def send_json(self, obj, code=200):
        d=json.dumps(obj).encode()
        self.send_response(code); self.send_header('Content-Type','application/json'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Access-Control-Allow-Methods','GET,POST,OPTIONS'); self.send_header('Content-Length',str(len(d))); self.end_headers(); self.wfile.write(d)
    def send_html(self):
        data=read_index()
        self.send_response(200); self.send_header('Content-Type','text/html; charset=utf-8'); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Cache-Control','no-cache'); self.send_header('Content-Length',str(len(data))); self.end_headers(); self.wfile.write(data)
    def do_GET(self):
        p=urlparse(self.path).path
        if p in ('/health','/ping'):
            self.send_response(200); self.send_header('Content-Length','2'); self.end_headers(); self.wfile.write(b'ok'); return
        if p=='/stats':
            total=cur.execute('SELECT SUM(balance), SUM(staked), COUNT(*) FROM balances').fetchone()
            txc=cur.execute('SELECT COUNT(*) FROM txs').fetchone()[0]
            return self.send_json({'chain':'EKI Chain MAINNET','chainId':CHAIN_ID,'block':block_now(),'total_supply':1000000,'circulating':float(total[0] or 0),'staked_total':float(total[1] or 0),'holders':int(total[2] or 0),'tx_count':txc,'rpc':'https://eki-chain.onrender.com/','status':'PRO'})
        if p.startswith('/balance/'):
            a=p.split('/balance/')[-1].lower()
            b,s=get_bal(a)
            return self.send_json({'address':a,'balance':b,'staked':s})
        if p=='/txs':
            rows=cur.execute('SELECT hash,frm,toaddr,amount,ts,block FROM txs ORDER BY block DESC LIMIT 50').fetchall()
            return self.send_json([{'hash':r[0],'from':r[1],'to':r[2],'amount':r[3],'timestamp':r[4],'block':r[5]} for r in rows])
        return self.send_html()
    def do_POST(self):
        l=int(self.headers.get('Content-Length',0)); body=self.rfile.read(l) if l else b'{}'
        try: req=json.loads(body or b'{}')
        except: req={}
        path=urlparse(self.path).path
        if path=='/claim':
            a=req.get('address','').lower()
            if not a.startswith('0x') or len(a)!=42: return self.send_json({'error':'endereco invalido'},400)
            if cur.execute('SELECT * FROM claimed WHERE addr=?',(a,)).fetchone(): return self.send_json({'error':'voce ja pegou, 1x por carteira'},400)
            b,s=get_bal(a); set_bal(a,b+10,s); cur.execute('INSERT INTO claimed(addr) VALUES(?)',(a,)); con.commit()
            return self.send_json({'ok':True,'balance':b+10,'msg':'Airdrop MAINNET PRO: 10 $EKI'})
        if path=='/stake':
            a=req.get('address','').lower(); amt=float(req.get('amount',0))
            b,s=get_bal(a)
            if b<amt: return self.send_json({'error':'saldo insuficiente'},400)
            set_bal(a,b-amt,s+amt); return self.send_json({'ok':True,'staked':s+amt,'apy':'12% ao ano'})
        if path=='/transfer':
            f=req.get('from','').lower(); t=req.get('to','').lower(); amt=float(req.get('amount',0))
            bf,sf=get_bal(f); bt,st=get_bal(t)
            if bf<amt: return self.send_json({'error':'saldo insuficiente'},400)
            blk=block_now(); h='0x'+hashlib.sha256(f"{f}{t}{amt}{time.time()}".encode()).hexdigest()
            set_bal(f,bf-amt,sf); set_bal(t,bt+amt,st)
            cur.execute('INSERT INTO txs VALUES(?,?,?,?,?,?)',(h,f,t,amt,int(time.time()),blk)); con.commit()
            return self.send_json({'ok':True,'hash':h,'from_balance':bf-amt,'to_balance':bt+amt})
        # --- JSON-RPC PRO - ISSO QUE A TRUST / CHAINLIST LE ---
        method=req.get('method',''); rid=req.get('id',1); params=req.get('params',[])
        blk=block_now(); res=None
        if method=='eth_chainId': res=CHAIN_HEX
        elif method=='net_version': res=str(CHAIN_ID)
        elif method=='web3_clientVersion': res='EKI Chain PRO/2.0.0'
        elif method=='net_listening': res=True
        elif method=='eth_blockNumber': res=to_hex(blk)
        elif method=='eth_gasPrice': res='0x0'
        elif method=='eth_maxPriorityFeePerGas': res='0x0'
        elif method=='eth_feeHistory': res={'baseFeePerGas':['0x0'],'gasUsedRatio':[0],'reward':[['0x0']],'oldestBlock':'0x0'}
        elif method=='eth_getBalance':
            addr=(params[0] if params else '').lower(); b,s=get_bal(addr); res=wei(b)
        elif method=='eth_getTransactionCount':
            addr=(params[0] if params else '').lower(); c=cur.execute('SELECT COUNT(*) FROM txs WHERE frm=?',(addr,)).fetchone()[0]; res=to_hex(c)
        elif method=='eth_estimateGas': res='0x5208'
        elif method=='eth_call':
            # Le balanceOf de token ERC20 - Trust usa isso pra mostrar token
            try:
                data=params[0].get('data','') if params else ''
                # balanceOf(address) = 0x70a08231 + 32 bytes endereco
                if data.startswith('0x70a08231'):
                    target='0x'+data[-40:].lower(); b,s=get_bal(target); res='0x'+format(int(b*10**18),'064x')
                else: res='0x'+'0'*64
            except: res='0x'+'0'*64
        elif method=='eth_getBlockByNumber':
            res={'number':to_hex(blk),'hash':'0x'+hashlib.sha256(str(blk).encode()).hexdigest(),'parentHash':'0x'+'0'*64,'timestamp':to_hex(int(time.time())),'miner':'0x'+'0'*40,'difficulty':'0x0','totalDifficulty':'0x0','size':'0x200','gasLimit':'0x1c9c380','gasUsed':'0x0','transactions':[],'uncles':[]}
        elif method=='eth_getTransactionByHash':
            h=params[0] if params else ''; r=cur.execute('SELECT frm,toaddr,amount,block FROM txs WHERE hash=?',(h,)).fetchone()
            res=None if not r else {'hash':h,'from':r[0],'to':r[1],'value':wei(r[2]),'blockNumber':to_hex(r[3]),'gasPrice':'0x0','gas':'0x5208'}
        elif method=='eth_getTransactionReceipt':
            h=params[0] if params else ''; r=cur.execute('SELECT frm,toaddr,block FROM txs WHERE hash=?',(h,)).fetchone()
            res=None if not r else {'transactionHash':h,'blockNumber':to_hex(r[2]),'status':'0x1','from':r[0],'to':r[1],'gasUsed':'0x5208'}
        elif method=='eth_getLogs': res=[]
        elif method=='eth_syncing': res=False
        else: res='0x0'
        return self.send_json({'jsonrpc':'2.0','id':rid,'result':res})
    def do_OPTIONS(self):
        self.send_response(200); self.send_header('Access-Control-Allow-Origin','*'); self.send_header('Access-Control-Allow-Headers','*'); self.send_header('Access-Control-Allow-Methods','GET,POST,OPTIONS'); self.end_headers()

PORT=int(os.environ.get('PORT',10000))
print(f"EKI PRO na porta {PORT} - Bloco {block_now()}", flush=True)
HTTPServer(('0.0.0.0',PORT), H).serve_forever()
