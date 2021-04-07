FROM ubuntu:20.04

COPY drb /

ENTRYPOINT ["/drb"]
