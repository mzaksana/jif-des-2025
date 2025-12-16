# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Copy source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 3000

# Run the server
CMD ["./server"]
