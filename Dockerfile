FROM alpine:latest

COPY build/rsec-net /bin/

ENTRYPOINT [ "/bin/rsec-net", "start-node", "-v"]
