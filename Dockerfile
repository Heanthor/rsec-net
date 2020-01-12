FROM alpine:latest

WORKDIR /app

ENTRYPOINT [ "build/linux/rsec-net", "start-node", "-v"]
