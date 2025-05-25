# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/share-a-ride ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/share-a-ride .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./share-a-ride"] 