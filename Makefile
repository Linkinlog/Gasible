MAIN := main.go
BINARY := gasible
PWD := $(shell pwd)

all: clean
all: $(BINARY)

$(BINARY): $(MAIN)
	go build -o $@ $^

clean:
	rm -f $(BINARY)

lint:
	golangci-lint run

fmt:
	gofmt -d -w $(PWD)

test:
	go test ./... -v

neat: lint fmt

check: neat test

.PHONY: all clean lint fmt test neat
