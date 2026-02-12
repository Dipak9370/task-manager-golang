APP_NAME=task-manager
MAIN_PATH=cmd/main.go
BIN_DIR=bin
BINARY=$(BIN_DIR)/$(APP_NAME)

.PHONY: all build run clean test swag setup

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BINARY) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./... -v

swag:
	swag fmt
	swag init -g $(MAIN_PATH)

setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

start: build
	$(BINARY)
