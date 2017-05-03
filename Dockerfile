FROM ubuntu

# build sysbench

RUN apt-get -qq update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y \
        ca-certificates \
        autoconf \
        libtool \
        git \
        pkg-config \
        vim \
    && apt-get clean -y \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* /tmp/* /var/tmp/*

# xxd was emancipated from vim in zesty

WORKDIR /root

RUN git clone https://github.com/akopytov/sysbench.git

RUN cd sysbench && \
    ./autogen.sh && \
    ./configure --without-mysql && \
    make && \
    cp src/sysbench /usr/local/bin/

ADD ./bin/beedrill-worker /usr/local/bin/beedrill-worker

ENTRYPOINT ["/usr/local/bin/beedrill-worker"]
