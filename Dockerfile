FROM ubuntu as sysbench

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

RUN cd sysbench \
    && ./autogen.sh \
    && ./configure --without-mysql \
    && make

# ==================================================================================

FROM ubuntu as filebench

RUN apt-get -qq update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y \
        ca-certificates \
        autoconf \
        automake \
        libtool \
        bison \
        flex \
        git \
        make \
    && apt-get clean -y \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /root

RUN git clone https://github.com/filebench/filebench.git

RUN cd filebench \
    && libtoolize \
    && aclocal \
    && autoheader \
    && automake --add-missing \
    && autoconf \
    && ./configure \
    && make

# ==================================================================================

FROM ubuntu

COPY --from=sysbench /root/sysbench/src/sysbench /usr/local/bin/
COPY --from=filebench /root/filebench/src/filebench /usr/local/bin/

COPY ./bin/beedrill-worker /usr/local/bin/beedrill-worker

ENTRYPOINT ["/usr/local/bin/beedrill-worker"]
