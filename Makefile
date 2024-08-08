BINARY_NAME := exec
SRC_DIR := .

all: build

build:
	@echo "Building the project..."
	@go build -o ./$(BINARY_NAME) $(SRC_DIR)

run: build
	@echo "Running the project..."
	@./$(BINARY_NAME)

deps:
	@echo "Installing dependencies..."
	@go mod tidy

deploy:
	@./deploy.sh

generate:
	@go run ./cmd/env.go ./.env
	
imports:
	@go run ./cmd/imports.go