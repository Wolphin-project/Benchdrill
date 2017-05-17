FROM ubuntu as sysbench

# build sysbench

RUN apt-get -qq update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y \
        ca-certificates \
        autoconf \
        libtool \
        libsasl2-dev \
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
    make

# ==================================================================================

FROM golang as beedrill

COPY . /go/src/git.rnd.alterway.fr/beedrill

WORKDIR /go/src/git.rnd.alterway.fr/beedrill

RUN go get github.com/urfave/cli/...
RUN go get github.com/RichardKnop/machinery/...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/beedrill ./cmd/beedrill.go

# ==================================================================================

FROM ubuntu

COPY --from=sysbench /root/sysbench/src/sysbench /usr/local/bin/
COPY --from=beedrill /go/src/git.rnd.alterway.fr/beedrill/bin/beedrill /usr/local/bin/

#ENTRYPOINT ["/usr/local/bin/beedrill", "worker"]
