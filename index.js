const express=require("express");
const app=express();
const path=require("path");
app.get("/chainlist.json",(req,res)=>res.sendFile(path.resolve("./chainlist.json")));
app.get("/explorer",(req,res)=>res.sendFile(path.resolve("./chainlist.json")));
app.get("/",(req,res)=>res.sendFile(path.resolve("./chainlist.json")));
app.post("/rpc",(req,res)=>res.json({jsonrpc:"2.0",id:req.body.id,result:"0x11165"}));
app.listen(process.env.PORT||10000);
