# Go LangChain-Style Agent for MCP Server

A pure Go implementation that provides LangChain-like orchestration capabilities for your MCP server, giving you intelligent tool coordination while staying in the Go ecosystem.

## 🚀 Features

### **Native Go Integration**
- ✅ **Pure Go**: No Python dependencies
- ✅ **Fast Performance**: Direct subprocess calls to MCP server
- ✅ **Type Safety**: Strong typing throughout
- ✅ **Easy Deployment**: Single binary

### **LangChain-Style Capabilities**
- 🧠 **Tool Orchestration**: Intelligent tool selection and chaining
- 💬 **Conversation Memory**: Context-aware conversations
- 🔄 **Multi-step Workflows**: Complex reasoning patterns
- 🎯 **Natural Language Interface**: GPT-powered command understanding

## 📦 Architecture

```
User Query
    ↓
OpenAI GPT (Tool Selection & Orchestration)
    ↓
Go Agent (Tool Execution)
    ↓
┌─────────────────┬─────────────────┐
│ MCPZipcodeTool  │ MCPClaudeTool   │
│       ↓         │       ↓         │
│ Go subprocess   │ Go subprocess   │
│       ↓         │       ↓         │
│ MCP Test Client │ MCP Test Client │
│       ↓         │       ↓         │
│   MCP Server    │   MCP Server    │
│       ↓         │       ↓         │
│   ViaCEP API    │   Claude API    │
└─────────────────┴─────────────────┘
```

## 🛠️ Setup

### 1. Install Dependencies
```powershell
cd go-agent
go mod tidy
```

### 2. Configure Environment
Make sure your root `.env` file contains:
```env
OPENAI_API_KEY=your_openai_key_here
CLAUDE_API_KEY=your_claude_key_here
```

### 3. Test the Tools
```powershell
go run test_simple.go
```

**Expected Output:**
```
🧪 Testing Go LangChain MCP Tools Integration
==================================================

📍 Testing Zipcode Tool:
✅ Success!
🔄 Initializing MCP connection...
✅ MCP connection initialized successfully
📤 Sending: {"id":2,"jsonrpc":"2.0","method":"tools/call"...}
📥 Response: {"id":2,"jsonrpc":"2.0","result":{"content":[{"text":"Your address is {...}

🤖 Testing Claude Tool:
✅ Success!
🔄 Initializing MCP connection...
✅ MCP connection initialized successfully
📤 Sending: {"id":3,"jsonrpc":"2.0","method":"tools/call"...}
📥 Response: {"id":3,"jsonrpc":"2.0","result":{"content":[{"text":"Claude says: 4"...}
```

### 4. Test All AI Providers
```powershell
go run test_all_ai.go
```

## 🎮 Usage

### **Basic Tool Testing**
```powershell
go run test_simple.go
```

### **Full Agent (Coming Soon)**
```powershell
go run main.go
```

## 🧠 What the Agent Can Do

### **Tool Wrapping**
The Go agent wraps your MCP server tools:

- **MCPZipcodeTool**: Brazilian address lookup
- **MCPClaudeTool**: Claude AI queries

### **Example Tool Usage**
```go
// Zipcode lookup
zipcodeTool := NewSimpleMCPZipcodeTool("../")
args := map[string]interface{}{"zipcode": "01310-100"}
result, err := zipcodeTool.Execute(args)

// Claude query
claudeTool := NewSimpleMCPClaudeTool("../")
args = map[string]interface{}{"question": "What is São Paulo?"}
result, err = claudeTool.Execute(args)
```

## 🔧 Current Implementation Status

### ✅ **Working Now**
- [x] Tool wrapper implementations
- [x] MCP server integration
- [x] Basic tool execution
- [x] Error handling
- [x] Clean JSON response parsing

### 🚧 **In Development**
- [ ] OpenAI integration for orchestration
- [ ] Conversation memory
- [ ] Tool chaining logic
- [ ] Interactive chat interface
- [ ] Complex workflow examples

## 📋 Benefits Over Python LangChain

### **Performance**
- ⚡ **Faster startup**: No Python interpreter overhead
- ⚡ **Lower memory**: Compiled binary vs interpreted
- ⚡ **Direct integration**: Same language as MCP server

### **Deployment**
- 📦 **Single binary**: Easy distribution
- 🔧 **No dependencies**: Self-contained executable
- 🚀 **Cloud native**: Perfect for containers

### **Development**
- 🔒 **Type safety**: Compile-time error checking
- 🧹 **Clean code**: Go's simplicity and clarity
- 🔄 **Fast iteration**: Quick compile-test cycles

## 🎯 Use Cases

### **Simple Workflows**
```
User: "Look up zipcode 01310-100"
→ Agent calls MCPZipcodeTool
→ Returns formatted address
```

### **Complex Workflows (Future)**
```
User: "Find the address for 01310-100 and tell me about that area"
→ Agent calls MCPZipcodeTool → Gets "Avenida Paulista"
→ Agent calls MCPClaudeTool → "Tell me about Avenida Paulista"
→ Agent combines results → Rich response
```

### **Conversational Context (Future)**
```
User: "Look up 01310-100"
Agent: [Returns address]

User: "What's special about that street?"
Agent: [Remembers context, asks Claude about Avenida Paulista]
```

## 🔍 Code Structure

### **Core Types**
```go
type Tool interface {
    Name() string
    Description() string
    Execute(args map[string]interface{}) (string, error)
}

type GoLangChainAgent struct {
    client    *openai.Client
    tools     map[string]Tool
    memory    *ConversationMemory
    systemMsg string
}
```

### **Tool Implementation**
```go
type MCPZipcodeTool struct {
    serverPath string
}

func (t *MCPZipcodeTool) Execute(args map[string]interface{}) (string, error) {
    // 1. Extract zipcode from args
    // 2. Call MCP test client subprocess
    // 3. Parse JSON response
    // 4. Return formatted result
}
```

## 🚨 Error Handling

The agent includes comprehensive error handling:
- Invalid tool arguments
- MCP server connection issues
- JSON parsing errors
- Subprocess execution failures

All errors are returned as structured messages for debugging.

## 🔄 Next Steps

1. **Complete OpenAI Integration**: Add the full agent orchestration
2. **Add Memory System**: Implement conversation context
3. **Tool Chaining**: Enable complex multi-step workflows
4. **Interactive Mode**: Command-line chat interface
5. **Web Interface**: HTTP API for the agent

## 🏗️ Why This Hybrid Approach Works

### **Keeps MCP Advantages**
- ✅ Fast, direct API calls
- ✅ MCP protocol compliance
- ✅ Lightweight core server
- ✅ Easy to maintain and debug

### **Adds LangChain-Style Power**
- 🧠 Intelligent tool orchestration (via OpenAI)
- 🔄 Multi-step reasoning
- 💬 Natural language interface
- 📚 Extensible tool system

### **Go-Specific Benefits**
- ⚡ Superior performance
- 🔒 Type safety and reliability
- 📦 Easy deployment and distribution
- 🧹 Clean, maintainable codebase

This Go agent gives you LangChain-like capabilities while maintaining the performance and simplicity of native Go! 🚀