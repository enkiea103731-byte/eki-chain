const express=require("express");
const app=express();
const CHAINLIST={"chainId":70101,"chainIdHex":"0x11165","chainlistFile":"https://eki-chain.onrender.com/chainlist.json","custody":{"BNB":"0x2Ba9927F45CA699E3FD6aC0Fd352A0569DBD0E88","BTC":"bc1qs03n0sft8p9aqzvfn8qg8smq7afh3g30ru53vt","EKI":"0x03664899007a53F086D545c34063E9CBfb29C9C2","ETH":"0x4DcC1EA8aaeD8fe1BdaF56694b37Adc3c2c973cf","SOL":"G3ERDDYZWu9gHNvQjzHep73ELMKTuXjKQDnr7MyRbenm"},"explorer":"https://eki-chain.onrender.com/explorer","network":"EKI MAINNET","rpc":"https://eki-chain.onrender.com/rpc","totalSupply":"1000000000000000","totalSupplyWei":"1000000000000000000000000000000"};
app.get("/chainlist.json",(req,res)=>res.json(CHAINLIST));
app.get("/explorer",(req,res)=>res.json(CHAINLIST));
app.get("/",(req,res)=>res.json(CHAINLIST));
app.post("/rpc",(req,res)=>res.json({jsonrpc:"2.0",id:req.body?.id||1,result:"0x11165"}));
app.listen(process.env.PORT||10000,()=>console.log("FIXED 70101 intacta"));
