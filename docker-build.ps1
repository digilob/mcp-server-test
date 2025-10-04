# Build and Run Multi-AI MCP Server with Docker
# PowerShell version for Windows

Write-Host "Building Multi-AI MCP Server Docker Images" -ForegroundColor Green
Write-Host "=" * 50 -ForegroundColor Blue

# Check if .env file exists
if (-not (Test-Path .env)) {
    Write-Host "WARNING: .env file not found!" -ForegroundColor Yellow
    Write-Host "Please copy .env.docker to .env and add your API keys:" -ForegroundColor White
    Write-Host "   Copy-Item .env.docker .env" -ForegroundColor Cyan
    Write-Host "   # Then edit .env with your API keys" -ForegroundColor Gray
    exit 1
}

# Build images
Write-Host "Building AI MCP Server..." -ForegroundColor Blue
docker build -t mcp-ai-server:latest .

Write-Host "Building File Operations Server..." -ForegroundColor Blue
docker build -t mcp-file-server:latest ./mcp-file-ops/

Write-Host "Build complete!" -ForegroundColor Green
Write-Host ""

Write-Host "Starting services with Docker Compose..." -ForegroundColor Blue
docker-compose up -d

Write-Host ""
Write-Host "Multi-AI MCP Server is now running!" -ForegroundColor Green
Write-Host "Services:" -ForegroundColor White
Write-Host "   AI MCP Server:    http://localhost:8080" -ForegroundColor Cyan
Write-Host "   File MCP Server:  http://localhost:8081" -ForegroundColor Cyan
Write-Host ""
Write-Host "Check status:" -ForegroundColor White
Write-Host "   docker-compose ps" -ForegroundColor Gray
Write-Host ""
Write-Host "View logs:" -ForegroundColor White
Write-Host "   docker-compose logs -f" -ForegroundColor Gray
Write-Host ""
Write-Host "Stop services:" -ForegroundColor White
Write-Host "   docker-compose down" -ForegroundColor Gray