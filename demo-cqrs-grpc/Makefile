.PHONY: proto build run clean demo

# Generate Go code from proto files
proto:
	@echo "Generating gRPC code from proto..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/blog.proto
	@echo "Done!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod tidy
	@echo "Done!"

# Build the server
build: deps
	@echo "Building server..."
	go build -o bin/server ./cmd/server
	@echo "Done! Binary at ./bin/server"

# Run the server
run: deps
	@echo "Starting CQRS Blog Service..."
	go run ./cmd/server

# Clean build artifacts and database
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f blog.db
	@echo "Done!"

# Run the demo script
demo:
	@chmod +x demo.sh
	@./demo.sh

# Help
help:
	@echo "CQRS + gRPC Blog Service Demo"
	@echo ""
	@echo "Commands:"
	@echo "  make proto  - Generate gRPC code from .proto files"
	@echo "  make deps   - Download Go dependencies"
	@echo "  make build  - Build the server binary"
	@echo "  make run    - Run the server"
	@echo "  make clean  - Remove build artifacts and database"
	@echo "  make demo   - Run the demo script"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Go 1.21+"
	@echo "  - protoc (Protocol Buffer compiler)"
	@echo "  - protoc-gen-go and protoc-gen-go-grpc"
	@echo "  - grpcurl (for testing)"
