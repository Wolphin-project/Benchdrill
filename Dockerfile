FROM ubuntu

COPY ./worker /worker

ENTRYPOINT ["/worker"]
