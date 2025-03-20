SHELL := /bin/bash

default: check

check: tidy fmt lint test

fmt:
	@echo "⏱️ formatting code now..."
	go fmt ./...
	@echo "✅ formatting finish"

build:
	@echo "⏱️ building..."
	go build -o ./bin/umex main.go
	@echo "✅ build finish"

test:
	@echo "⏱️ running tests now... "
	go test -race --parallel=4 -timeout 30s -cover $(ARGS)
	@echo "✅ passing all tests."

lint:
	@echo "⏱️ running linting now..."
	golangci-lint run $(ARGS)
	@echo "✅ passing linting..."

tidy:
	@echo "⏱️ go mod tidy now..."
	go mod tidy
	@echo "✅ finishing tidy..."