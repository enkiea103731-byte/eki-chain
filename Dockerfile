FROM ghcr.io/foundry-rs/foundry:latest
EXPOSE 10000
ENTRYPOINT ["anvil","--host","0.0.0.0","--port","10000","--chain-id","70101"]
