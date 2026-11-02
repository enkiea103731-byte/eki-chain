const express = require('express');
const app = express();
app.use(express.json());

const CHAIN_DATA = {
  "chainId": 70101,
  "chainIdHex": "0x111D5",
  "name": "EKI 70101",
  "networkId": 70101,
  "rpc": ["https://eki-chain.onrender.com"],
  "explorer": "https://eki-chain.onrender.com",
  "infoURL": "https://eki-chain.onrender.com",
  "status": "Online",
  "custody": {
    "BNB": "0x2Ba9927F45CA699E3FD6aC0fD352A0569DBDE88",
    "BTC": "bc1q03n0sTf8p9aqzvfn8qg8smq7afh3g30ru53vt",
    "EKI": "0x03664899007a53f086d545c34063e9cbfb29c9c2",
    "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
    "SOL": "G3ERDDYZhU9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm"
  }
};

app.get('/', (req, res) => {
  res.json(CHAIN_DATA);
});

app.post('/', (req, res) => {
  const method = req.body?.method;
  if (method === 'eth_chainId') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x111D5" });
  }
  if (method === 'net_version') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "70101" });
  }
  if (method === 'eth_blockNumber') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x1" });
  }
  return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x0" });
});

const PORT = process.env.PORT || 10000;
app.listen(PORT, () => console.log('EKI 70101 Online na porta ' + PORT));
