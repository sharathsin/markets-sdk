.PHONY: all build run-crypto run-stock test clean

APP_NAME=markets

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/markets/main.go

run-crypto: build
	@./bin/$(APP_NAME) -provider crypto -symbol bitcoin

run-stock: build
	@./bin/$(APP_NAME) -provider stock -symbol AAPL

test:
	@go test ./... -v

clean:
	@rm -rf bin
	@echo "Cleaned"
