BINARIES=bin/worker bin/send bin/worker_scratch

all: $(BINARIES)

bin/worker:
	go build -o bin/worker ./cmd/worker.go

bin/send:
	go build -o bin/send ./cmd/send.go

bin/worker_scratch:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/worker_scratch ./cmd/worker.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
