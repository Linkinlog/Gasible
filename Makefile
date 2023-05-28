MAIN := main.go
BINARY := gas
PWD := $(shell pwd)

all: clean
all: $(BINARY)

$(BINARY): $(MAIN)
	go build -o $@ $^

clean:
	rm -f $(BINARY)

lint:
	command golangci-lint run

fmt:
	command gofmt -d -w $(PWD)

test:
	command go test ./... -v

neat: lint fmt

check: neat test

.PHONY: all clean lint fmt test neat
