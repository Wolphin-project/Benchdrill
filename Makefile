BINARIES=bin/beedrill bin/beedrill-worker

all: $(BINARIES)

bin/beedrill:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/beedrill ./cmd/beedrill.go

bin/beedrill-worker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/beedrill-worker ./cmd/beedrill.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
