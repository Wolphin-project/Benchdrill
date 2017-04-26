FROM ubuntu

ADD ./bin/beedrill-worker /beedrill-worker

RUN ["/bin/bash", "-c", "apt-get update && apt-get install -y curl"]

RUN ["/bin/bash", "-c", "curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.deb.sh | bash"]

RUN apt-get update && apt-get install -y sysbench

ENTRYPOINT ["/beedrill-worker"]
