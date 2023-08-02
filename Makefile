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
	golangci-lint --color=always --timeout=30m --exclude-use-default=false --print-linter-name -E=govet -E=asciicheck -E=bodyclose -E=dupl -E=gocognit -E=gocyclo -E=goerr113 -E=unused -E=gomnd -E=gosec -E=misspell -E=nestif -E=rowserrcheck -E=unconvert -E=unparam -E=whitespace -E=goconst run

fmt:
	gofmt -d -w $(PWD)

test:
	go test ./... -v

neat: lint fmt

check: neat test

.PHONY: all clean lint fmt test neat
