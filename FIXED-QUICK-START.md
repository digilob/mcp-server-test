# 🚨 CORRECTED Quick Start Guide

## The Problem
The documentation was incorrect. The "integrated client" (`cmd/mcp-test-client`) **does NOT work** - it has timeout issues.

## The Solution
Use the **manual two-terminal approach** which actually works:

### Working Method ✅

**Terminal 1: Start Server**
```powershell
# Make sure you're in the PROJECT ROOT (not in any cmd/ subdirectory)
cd C:\Users\lobth\projects\mcp-server-test
go run .\main.go
```

**Terminal 2: Use Client**
```powershell
# Start from project root, then go to test-client
cd C:\Users\lobth\projects\mcp-server-test
cd cmd\test-client
go run .\main.go list
go run .\main.go zipcode 01310-100
```

### Important: Directory Check ⚠️
Make sure you're in the **project root** for Terminal 1, not in any `cmd/` subdirectory:
```powershell
# ✅ CORRECT - You should be here for Terminal 1:
PS C:\Users\lobth\projects\mcp-server-test> go run .\main.go

# ❌ WRONG - Don't run server from here:
PS C:\Users\lobth\projects\mcp-server-test\cmd\mcp-test-client> go run .\main.go
```

### What You'll See (Working)
```
📤 Sending: {"id":1,"jsonrpc":"2.0","method":"tools/list","params":{}}
📥 Response: {"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"zipcode","description":"Find an address by zip code"}...]}}
```

### Broken Method ❌ (Don't Use)
```powershell
cd cmd\mcp-test-client
go run .\main.go claude "test"
# Results in: ❌ Timeout waiting for response
```

## Why This Happened
The integrated client tries to start its own server process but has initialization and timing issues. The manual approach works because:
1. Server starts cleanly
2. Client connects to already-running server
3. No timing conflicts

## Updated Documentation
The README.md has been corrected to show the **working method first** and mark the broken integrated client as problematic.

Use this approach and you'll get proper MCP responses instead of timeouts!

## 🔧 Troubleshooting

### Problem: "Usage: go run main.go <command>"
**Cause**: You're trying to run the server from the wrong directory
**Solution**: Navigate to project root first
```powershell
# Check where you are
pwd

# If you see cmd\mcp-test-client or cmd\test-client, go back:
cd ..\..

# Now you should be in the project root:
PS C:\Users\lobth\projects\mcp-server-test> go run .\main.go
```

### Problem: "Timeout waiting for response"  
**Cause**: Using the broken integrated client
**Solution**: Use the manual two-terminal method above

### Problem: Server won't start
**Cause**: Missing dependencies or wrong directory
**Solution**: 
```powershell
cd C:\Users\lobth\projects\mcp-server-test
go mod tidy
go run .\main.go
```