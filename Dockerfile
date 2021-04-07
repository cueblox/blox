FROM ubuntu:20.04

COPY blox /

ENTRYPOINT ["/blox"]
