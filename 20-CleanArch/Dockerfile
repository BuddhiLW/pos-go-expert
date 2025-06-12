# syntax=docker/dockerfile:1.4
FROM golang:latest AS builder

WORKDIR /app

# Install git for go mod download
RUN apt-get update && apt-get install -y git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ordersystem ./cmd/ordersystem

FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/ordersystem .

# Expose ports
EXPOSE 8000 50051 9002

# Run the application
CMD ["./ordersystem"] 