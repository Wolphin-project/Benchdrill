BINARIES=bin/worker bin/send

all: $(BINARIES)

bin/worker:
	go build -o bin/worker ./worker/worker.go

bin/send:
	go build -o bin/send ./send/send.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
