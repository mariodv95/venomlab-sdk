.PHONY: build run clean test help

# Variables
GO_OUT=proto/build
GO_GRPC_OUT=proto/build
PROTO_FILES=proto/*.proto

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Venomlab SDK Proto Files
	@echo "Building Proto..."
	protoc -I proto --go_out=$(GO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(GO_GRPC_OUT) --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	@echo "Proto built"

