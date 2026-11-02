const http = require('http');
const DATA = {"chainId":70101,"chainIdHex":"0x111D5","name":"EKI 70101","rpc":["https://eki-chain.onrender.com"],"explorer":"https://eki-chain.onrender.com","infoURL":"https://eki-chain.onrender.com","networkId":70101,"status":"Online"};
const server = http.createServer((req,res)=>{
  let b=''; req.on('data',c=>b+=c); req.on('end',()=>{
    res.setHeader('Content-Type','application/json');
    res.setHeader('Access-Control-Allow-Origin','*');
    if(req.method==='POST'){
      try{
        const j=JSON.parse(b||'{}');
        if(j.method==='eth_chainId') return res.end(JSON.stringify({jsonrpc:"2.0",id:j.id||1,result:"0x111D5"}));
        if(j.method==='net_version') return res.end(JSON.stringify({jsonrpc:"2.0",id:j.id||1,result:"70101"}));
      }catch(e){}
      return res.end(JSON.stringify({jsonrpc:"2.0",id:1,result:"0x111D5"}));
    }
    res.end(JSON.stringify(DATA));
  });
});
server.listen(process.env.PORT||10000);
