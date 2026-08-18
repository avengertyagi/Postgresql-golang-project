BINARY_DIR := bin
API_BIN    := $(BINARY_DIR)/api

.PHONY: build run dev start seed keygen tidy fmt vet test clean help

## build: compile the API binary into bin/api
build:
	@mkdir -p $(BINARY_DIR)
	go build -o $(API_BIN) ./cmd/api

## run: go build + run in one go, mirrors `go build ./cmd/api/main.go && ./main`
run: build
	./$(API_BIN)

## dev: run without building a binary first, mirrors `go run ./cmd/api/main.go`
dev:
	go run ./cmd/api

## start: run an already-built binary, mirrors `./main` 
start:
	./$(API_BIN)

## seed: mirrors `go run ./cmd/seeder/main.go`
seed:
	go run ./cmd/seeder

## keygen: mirrors `go run ./cmd/keygen/main.go`
keygen:
	go run ./cmd/keygen

## tidy: sync go.mod/go.sum
tidy:
	go mod tidy

## fmt: format all Go source files
fmt:
	go fmt ./...

## vet: static analysis, catches suspicious code before it becomes a bug
vet:
	go vet ./...

## test: run all tests
test:
	go test ./... -v

## clean: remove built binaries and local log file
clean:
	rm -rf $(BINARY_DIR) gin.log

## help: list all available commands with their description
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  make /'