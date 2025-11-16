.PHONY: lint lint-fix test build deps

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

fmt:
	go fmt ./...

imports:
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -l -w .

lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

lint-fix:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --fix ./...

test:
	go test -v -race ./...

check: fmt lint test

check-all: deps fmt imports lint test

build:
	go build -o bin/app ./cmd/main.go

run:
	go run ./cmd/main.go


clean:
	rm -rf bin/