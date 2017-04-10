FROM scratch

ADD ./bin/worker /worker

ENTRYPOINT ["/worker"]
