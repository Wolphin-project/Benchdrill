BINARIES=bin/worker bin/send

all: $(BINARIES)

bin/worker:
	go build -o bin/worker ./cmd/worker.go

bin/send:
	go build -o bin/send ./cmd/send.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
