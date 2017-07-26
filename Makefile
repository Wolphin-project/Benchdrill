BINARIES=bin/benchdrill

all: $(BINARIES)

bin/benchdrill:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/benchdrill ./cmd/benchdrill.go

clean:
	rm -rf ${BINARIES}

.PHONY: clean
