FROM scratch

ADD ./bin/beedrill-worker /beedrill-worker

ENTRYPOINT ["/beedrill-worker"]
