FROM ubuntu

COPY ./bin/worker /worker

ENTRYPOINT ["/worker"]
