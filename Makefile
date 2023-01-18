lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v
fmt:
	gofmt -w ./

.PHONY: lint fmt
