# Build stage
FROM golang:1.23-alpine AS builder

# Install protoc and plugins
RUN apk add --no-cache protobuf protobuf-dev git make gcc musl-dev

WORKDIR /app

# Install Go protoc plugins (specific versions compatible with Go 1.23)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# Copy source code
COPY . .

# Download dependencies
RUN go mod tidy

# Generate protobuf code
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/blog.proto

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.18

RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose gRPC port
EXPOSE 50051

# Run the server
CMD ["./server"]
