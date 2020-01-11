FROM alpine:latest

WORKDIR /app

ENTRYPOINT [ "build/rsec-net", "start-node"]
