# 🐳 Docker Deployment Guide

Complete guide for running the Multi-AI MCP Server using Docker.

## 🚀 Quick Start

### 1. **Prerequisites**
- Docker Desktop installed
- Docker Compose available
- API keys for AI providers

### 2. **Environment Setup**
```powershell
# Copy environment template
Copy-Item .env.docker .env

# Edit .env with your API keys
notepad .env
```

### 3. **Build and Run**
```powershell
# Using PowerShell script (Windows)
.\docker-build.ps1

# Or manually with Docker Compose
docker-compose up -d
```

## 📦 Container Architecture

```
┌─────────────────────────────────────────┐
│              Docker Host                │
├─────────────────┬───────────────────────┤
│  mcp-ai-server  │   mcp-file-server     │
│  Port: 8080     │   Port: 8081          │
│  ┌─────────────┐│   ┌─────────────────┐ │
│  │5 AI Providers││   │File Operations  │ │
│  │- Claude     ││   │- Read/Write     │ │
│  │- OpenAI     ││   │- Search         │ │
│  │- Gemini     ││   │- List           │ │
│  │- Mistral    ││   │- Info           │ │
│  │- HuggingFace││   │                 │ │
│  │- Zipcode    ││   │                 │ │
│  └─────────────┘│   └─────────────────┘ │
└─────────────────┴───────────────────────┘
```

## 🐳 Available Images

### **Main AI MCP Server**
- **Image**: `mcp-ai-server:latest`
- **Port**: 8080
- **Features**: 5 AI providers + zipcode lookup
- **Protocol**: JSON-RPC over stdio

### **File Operations Server**
- **Image**: `mcp-file-server:latest`  
- **Port**: 8081
- **Features**: File read/write/search operations
- **Volume**: `/data` for persistent file storage

## 🔧 Configuration

### **Environment Variables**
| Variable | Description | Required |
|----------|-------------|----------|
| `CLAUDE_API_KEY` | Anthropic Claude API key | ✅ |
| `OPENAI_API_KEY` | OpenAI GPT API key | ✅ |
| `GEMINI_API_KEY` | Google Gemini API key | ✅ |
| `MISTRAL_API_KEY` | Mistral AI API key | ✅ |
| `HUGGINGFACEHUB_API_TOKEN` | Hugging Face token | ✅ |

### **Ports**
- **8080**: Main AI MCP Server
- **8081**: File Operations Server

### **Volumes**
- `./data:/data` - File operations working directory
- `./temp:/tmp` - Temporary files

## 🚀 Usage Examples

### **Start Services**
```powershell
# Build and start
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### **Test the Services**
```powershell
# Connect to AI server container
docker exec -it mcp-ai-server /bin/sh

# Connect to file server container  
docker exec -it mcp-file-server /bin/sh
```

### **Stop Services**
```powershell
# Stop containers
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 🔍 Health Checks

Both containers include health checks:

```powershell
# Check container health
docker ps

# View health check logs
docker inspect mcp-ai-server | grep Health -A 10
```

## 📊 Monitoring

### **View Live Logs**
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f mcp-ai-server
docker-compose logs -f mcp-file-server
```

### **Resource Usage**
```powershell
# Container stats
docker stats

# Specific container
docker stats mcp-ai-server
```

## 🔧 Development Mode

### **Development Compose Override**
Create `docker-compose.override.yml`:
```yaml
version: '3.8'
services:
  mcp-ai-server:
    volumes:
      - .:/app
    command: ["go", "run", "main.go"]
    
  mcp-file-server:
    volumes:
      - ./mcp-file-ops:/app
    command: ["go", "run", "main.go"]
```

### **Hot Reload Development**
```powershell
# Development mode with volume mounts
docker-compose -f docker-compose.yml -f docker-compose.override.yml up
```

## 🚀 Production Deployment

### **Docker Swarm**
```powershell
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml mcp-stack
```

### **Kubernetes**
```powershell
# Generate Kubernetes manifests
docker-compose config > mcp-k8s.yml

# Apply to cluster
kubectl apply -f mcp-k8s.yml
```

## 🛡️ Security Best Practices

### **API Key Management**
- ✅ Use `.env` files (not committed to Git)
- ✅ Consider Docker secrets for production
- ✅ Rotate API keys regularly

### **Network Security**
- ✅ Use Docker networks for service communication
- ✅ Expose only necessary ports
- ✅ Consider reverse proxy for production

### **Container Security**
- ✅ Run as non-root user (implemented)
- ✅ Use minimal Alpine base images
- ✅ Regular security updates

## 🔧 Troubleshooting

### **Common Issues**

#### **Container Won't Start**
```powershell
# Check logs
docker-compose logs mcp-ai-server

# Check configuration
docker-compose config
```

#### **API Key Errors**
```powershell
# Verify environment
docker exec mcp-ai-server env | grep API

# Test .env file
Get-Content .env
```

#### **Network Issues**
```powershell
# Check network
docker network ls

# Inspect network
docker network inspect mcp-server-test_mcp-network
```

## 📈 Scaling

### **Horizontal Scaling**
```yaml
# docker-compose.yml
services:
  mcp-ai-server:
    deploy:
      replicas: 3
```

### **Load Balancing**
Add nginx or traefik for load balancing multiple instances.

## 🎯 Next Steps

1. **Production Deployment**: Deploy to cloud providers
2. **Monitoring**: Add Prometheus/Grafana
3. **CI/CD**: Automate builds and deployments
4. **Web UI**: Add browser interface
5. **API Gateway**: Add rate limiting and auth

Your Multi-AI MCP Server is now fully containerized and ready for any deployment scenario! 🚀