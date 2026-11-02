const express = require('express');
const app = express();
app.use(express.json());

const CHAIN_INFO = {
  "cadeia": "EKI",
  "chainId": 70101,
  "chainSlug": "eki",
  "custodia": {
    "BNB": "0x2Bb9027F45CA699E3FD6aC0fD352A0569DBDE88",
    "BTC": "bc1q03n0sTf8p9aqzvfn8qg8smq7afh3g30ru53vt",
    "ETH": "0xD4DCfCEAa8ae8f8e1eBdAf56694b37Adc3c2c973cf",
    "SOL": "G3ERDDYZhU9gHNvQjzHep73ELMKTuXjkQDnr7MyRbenm"
  },
  "exploradores": [],
  "taxa": 0.003,
  "freeWallet": "EKI_3le3da3e06a0af3",
  "infoURL": "https://eki-chain.onrender.com",
  "nome": "EKI Chain",
  "moeda_nativa": {
    "decimals": 18,
    "nome": "EKI",
    "simbolo": "EKI"
  },
  "networkId": 70101,
  "rpc": ["https://eki-chain.onrender.com"],
  "nome_abreviado": "eki",
  "status": "Online"
};

app.get('/', (req, res) => {
  res.json(CHAIN_INFO);
});

// Responde eth_chainId com 70101 (0x111D5) - FIX DEFINITIVO
app.post('/', (req, res) => {
  const method = req.body.method;
  if (method === 'eth_chainId') {
    return res.json({ jsonrpc: "2.0", id: req.body.id, result: "0x111D5" });
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
app.listen(PORT, () => console.log('Cadeia EKI 70101 Online em ' + PORT));
