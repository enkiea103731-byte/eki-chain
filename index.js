const CHAINLIST = {
  "chainId": 70101,
  "chainIdHex": "0x111D5",
  "chainlistFile": "https://eki-chain.onrender.com/chainlist.json",
  "custody": {
    "BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
    "BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
    "EKI": "0x03664899007a53f086d545c34063e9cbfb29c9c2",
    "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
    "SOL": "G3ERDDYZ Wu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm"
  },
  "explorer": "https://eki-chain.onrender.com/explorer",
  "network": "EKI MAINNET",
  "rpc": "https://eki-chain.onrender.com/rpc",
  "totalSupply": "1000000000000000",
  "totalSupplyWei": "1000000000000000000000000000000"
};

const http = require('http');
http.createServer((req,res)=>{
  if(req.url.includes('chainId')){
    res.writeHead(200,{'Content-Type':'application/json'});
    res.end(JSON.stringify({jsonrpc:"2.0",id:1,result:"0x111D5"}));
  } else {
    res.writeHead(200,{'Content-Type':'application/json'});
    res.end(JSON.stringify(CHAINLIST));
  }
}).listen(process.env.PORT||10000);
