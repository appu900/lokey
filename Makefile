BINARY := lokey
BIN_DIR := bin
PORT := 1111

.DEFAULT_GOAL := help


build:
	go build -o $(BIN_DIR)/$(BINARY) .

## run: run the server on port 1111 without building a binary
run:
	go run .

## start: build, then run the compiled binary
start: build
	./$(BIN_DIR)/$(BINARY)

## test: run all tests
test:
	go test ./...

## test-v: run all tests verbosely, skipping the result cache
test-v:
	go test ./... -v -count=1

## cover: run tests and open the coverage report in a browser
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

## fmt: format all Go source
fmt:
	go fmt ./...

## vet: report suspicious constructs
vet:
	go vet ./...

## tidy: sync go.mod with the imports actually used
tidy:
	go mod tidy

## check: fmt, vet and test - run this before committing
check: fmt vet test

## clean: remove build output and coverage data
clean:
	rm -rf $(BIN_DIR) coverage.out

## help: list the available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  make /'

.PHONY: build run start test test-v cover fmt vet tidy check clean help
