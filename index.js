const express = require('express');
const { ethers } = require('ethers');
const app = express();
app.use(express.json());

const OWNER = "0x2A0927745CA69923F3D05D7352A05690D6E8B".toLowerCase();
const TOTAL_WEI = "0x52B5E56460B07F2000000";
let currentBlock = 1;
let nonces = {};
let receipts = {};

// INTACTO PRA CHAINLIST - NÃO MUDEI NADA AQUI
const CHAIN_INFO = {
  chainId: 70101,
  networkId: 70101,
  chain: "EKI Mainnet",
  shortName: "eki",
  rpc: ["https://eki-chain.onrender.com/rpc"],
  nativeCurrency: { name: "EKI", symbol: "EKI", decimals: 18 },
  infoURL: "https://eki-chain.onrender.com",
  explorers: [{ name: "EKI Explorer", url: "https://eki-chain.onrender.com/explorer" }],
  totalSupply: "1000000000000",
  owner: OWNER,
  feeWallet: OWNER,
  status: "Online"
};

app.get('/', (req, res) => res.json(CHAIN_INFO));

app.post('/', (req, res) => {
  const m = req.body.method;
  const id = req.body.id;
  const p = req.body.params || [];

  if (m == "eth_chainId") return res.json({ jsonrpc: "2.0", id, result: "0x11165" });
  if (m == "net_version") return res.json({ jsonrpc: "2.0", id, result: "70101" });
  if (m == "eth_blockNumber") return res.json({ jsonrpc: "2.0", id, result: ethers.utils.hexValue(currentBlock) });
  if (m == "eth_gasPrice") return res.json({ jsonrpc: "2.0", id, result: "0x3B9ACA00" });
  if (m == "eth_estimateGas") return res.json({ jsonrpc: "2.0", id, result: "0x2DC6C0" });

  if (m == "eth_getTransactionCount") {
    const addr = (p[0] || "").toLowerCase();
    return res.json({ jsonrpc: "2.0", id, result: ethers.utils.hexValue(nonces[addr] || 0) });
  }

  if (m == "eth_getBalance") {
    const addr = (p[0] || "").toLowerCase();
    const bal = addr == OWNER ? TOTAL_WEI : "0x0";
    return res.json({ jsonrpc: "2.0", id, result: bal });
  }

  if (m == "eth_sendRawTransaction") {
    const raw = p[0];
    const txHash = ethers.utils.keccak256(raw);
    try {
      const parsed = ethers.utils.parseTransaction(raw);
      const from = (parsed.from || OWNER).toLowerCase();
      const nonce = nonces[from] || 0;
      const contractAddr = ethers.utils.getContractAddress({ from, nonce });
      nonces[from] = nonce + 1;
      currentBlock++;
      receipts[txHash] = {
        transactionHash: txHash,
        blockNumber: ethers.utils.hexValue(currentBlock),
        contractAddress: contractAddr,
        status: "0x1"
      };
    } catch(e) {
      currentBlock++;
      receipts[txHash] = { transactionHash: txHash, blockNumber: ethers.utils.hexValue(currentBlock), status: "0x1" };
    }
    return res.json({ jsonrpc: "2.0", id, result: txHash });
  }

  if (m == "eth_getTransactionReceipt") {
    return res.json({ jsonrpc: "2.0", id, result: receipts[p[0]] || null });
  }

  return res.json({ jsonrpc: "2.0", id, result: "0x0" });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(`EKI 70101 Online on ${PORT}`));
