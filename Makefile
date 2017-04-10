BINARIES=bin/worker bin/send

all: $(BINARIES)

bin/worker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/worker ./cmd/worker.go

bin/send:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/send ./cmd/send.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
