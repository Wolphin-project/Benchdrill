BINARIES=bin/beedrill

all: $(BINARIES)

bin/beedrill:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/beedrill ./cmd/beedrill.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
