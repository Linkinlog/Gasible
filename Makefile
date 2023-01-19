pwd := $(shell pwd)
lint:
	command golangci-lint run
fmt:
	command gofmt -d -w $(pwd)
neat: lint fmt

.PHONY: lint fmt neat
