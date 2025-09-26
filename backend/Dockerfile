# =========================
# Build stage
# =========================
FROM golang:1.25.1 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source code
COPY . .

# Build the Go binary
# main.go is now directly under backend/, so adjust path if Dockerfile is at repo root
RUN go build -o cyber-service main.go

# =========================
# Final stage
# =========================
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/cyber-service .

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./cyber-service"]