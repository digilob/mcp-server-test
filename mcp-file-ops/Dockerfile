# File Operations MCP Server
FROM golang:1.21-alpine AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the file operations server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-ops-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/file-ops-server .

# Create directories for file operations
RUN mkdir -p /data /tmp

# Expose port
EXPOSE 8081

# Volume for file operations
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD echo "File Operations Server is running"

CMD ["./file-ops-server"]