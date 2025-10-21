# Use a newer Go version (1.25+)
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the app (static binary)
RUN go build -o shortener cmd/api/main.go

# ---- Runtime Stage ----
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/shortener .

# Expose the port
EXPOSE 8080

# Run the app
CMD ["./shortener"]

