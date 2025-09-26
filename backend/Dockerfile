FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cyber-service main.go

EXPOSE 8080

CMD ["./cyber-service"]
# Use official Go image
FROM golang:1.25.1

# Set working directory to your app root
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the server from cmd/server
RUN go build -o cyber-go ./cmd/server

# Expose port
EXPOSE 8080

# Run the server
CMD ["./cyber-go"]