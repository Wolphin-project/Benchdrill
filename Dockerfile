FROM ubuntu

RUN apt-get update && apt-get install -y curl
RUN curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.deb.sh | bash
RUN apt-get update && apt-get install -y sysbench

COPY ./bin/beedrill-worker /beedrill-worker

ENTRYPOINT ["/beedrill-worker"]
