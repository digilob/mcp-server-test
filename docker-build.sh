#!/bin/bash

# Build and Run Multi-AI MCP Server with Docker

echo "🚀 Building Multi-AI MCP Server Docker Images"
echo "=============================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "Please copy .env.docker to .env and add your API keys:"
    echo "   cp .env.docker .env"
    echo "   # Then edit .env with your API keys"
    exit 1
fi

# Build images
echo "📦 Building AI MCP Server..."
docker build -t mcp-ai-server:latest .

echo "📦 Building File Operations Server..."
docker build -t mcp-file-server:latest ./mcp-file-ops/

echo "🎉 Build complete!"
echo ""
echo "🚀 Starting services with Docker Compose..."
docker-compose up -d

echo ""
echo "✅ Multi-AI MCP Server is now running!"
echo "📊 Services:"
echo "   • AI MCP Server:    http://localhost:8080"
echo "   • File MCP Server:  http://localhost:8081"
echo ""
echo "🔍 Check status:"
echo "   docker-compose ps"
echo ""
echo "📝 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down"