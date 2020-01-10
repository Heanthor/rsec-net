FROM alpine:latest

WORKDIR /app

# ENTRYPOINT [ "/bin/rsec-net", "start-node", "-v"]
ENTRYPOINT [ "build/rsec-net", "start-node", "-v", "--profile", "--profileMode=goroutine"]
