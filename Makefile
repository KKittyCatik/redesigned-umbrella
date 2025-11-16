.PHONY: lint lint-fix test build deps

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

fmt:
	go fmt ./...

imports:
	goimports -l -w .

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test -v -race ./...

check: fmt lint test

check-all: deps fmt imports lint test

build:
	go build -o ./bin/redesigned-umbrella ./cmd/pr-reviewer/main.go

run:
	./bin/redesigned-umbrella

clean:
	rm -rf bin/