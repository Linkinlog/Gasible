pwd := $(shell pwd)
lint:
	command golangci-lint run
fmt:
	command gofmt -d -w $(pwd)
test:
	command go test ./...
neat: lint fmt
check: neat test

.PHONY: lint fmt test neat
