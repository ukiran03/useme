# Variables
BINARY_NAME=app

# .PHONY tells Make that these targets aren't actual files
.PHONY: build run clean help

build: ## Build the binary
	mkdir -p build
	go build -o build/${BINARY_NAME} .

run: build ## Build and run the app (don't use yet)
	./build/${BINARY_NAME} $(ARGS)

clean: ## Remove the binary and build folder
	rm -rf build/
	go clean

help: ## Show this help focus
	@echo "Usage: make [target]"
	@echo "targets: "
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
