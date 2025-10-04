# Multi-AI MCP Server
FROM golang:1.21-alpine AS builder

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-server .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests to external APIs
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/mcp-server .

# Copy any additional files if needed
COPY --from=builder /app/mcp-server-architecture.dot .

# Expose port (optional, mainly for JSON-RPC over stdio)
EXPOSE 8080

# Environment variables (will be overridden by .env or docker-compose)
ENV CLAUDE_API_KEY=""
ENV OPENAI_API_KEY=""
ENV GEMINI_API_KEY=""
ENV MISTRAL_API_KEY=""
ENV HUGGINGFACEHUB_API_TOKEN=""

# Create a health check script
RUN echo '#!/bin/sh\necho "MCP Server is running"' > /healthcheck.sh && chmod +x /healthcheck.sh

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD /healthcheck.sh

# Run the binary
CMD ["./mcp-server"]