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

FROM golang as beedrill

RUN apt-get -qq update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y \
        ca-certificates \
        libsasl2-dev \
    && apt-get clean -y \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY . /go/src/git.rnd.alterway.fr/Wolphin-project/beedrill

WORKDIR /go/src/git.rnd.alterway.fr/Wolphin-project/beedrill

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/beedrill ./cmd/beedrill.go

# ==================================================================================

FROM ubuntu

COPY --from=sysbench /root/sysbench/src/sysbench /usr/local/bin/
COPY --from=filebench /root/filebench/filebench /usr/local/bin/
COPY --from=beedrill /go/src/git.rnd.alterway.fr/Wolphin-project/beedrill/bin/beedrill /usr/local/bin/

WORKDIR /root

COPY config_beedrill.yml /root/
COPY config_beedrill-worker.yml /root/

ENTRYPOINT ["/usr/local/bin/beedrill"]
