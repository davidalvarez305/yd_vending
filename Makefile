BINARY_NAME := budgeting
BUILD_DIR := build
SRC_DIR := .

all: build

build:
	@echo "Building the project..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

run: build
	@echo "Running the project..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

deps:
	@echo "Installing dependencies..."
	@go mod tidy

deploy:
	@./deploy.sh