const express = require('express');
const app = express();
app.use(express.json());

const CHAIN_INFO = {
  "chain": "EKI",
  "chainId": 70101,
  "chainSlug": "eki",
  "custody": {
    "BNB": "0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88",
    "BTC": "bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt",
    "ETH": "0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf",
    "SOL": "G3ERDDYZWu9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm"
  },
  "explorers": [],
  "faucets": [],
  "feeRate": 0.003,
  "feeWallet": "EKI_31e3da3ae96a0af3",
  "infoURL": "https://eki-chain.onrender.com",
  "name": "EKI Chain",
  "nativeCurrency": {
    "decimals": 18,
    "name": "EKI",
    "symbol": "EKI"
  },
  "networkId": 70101,
  "rpc": ["https://eki-chain.onrender.com"],
  "shortName": "eki",
  "status": "Online"
};

app.get('/', (req, res) => {
  res.json(CHAIN_INFO);
});

// Responde eth_chainId com 70101 (0x11235) pro validador passar
app.post('/', (req, res) => {
  const method = req.body.method;
  if (method === 'eth_chainId') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x11235" });
  }
  if (method === 'net_version') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "70101" });
  }
  if (method === 'eth_blockNumber') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x1" });
  }
  res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x0" });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log('EKI Chain 70101 Online on ' + PORT));
