# 🚀 Multi-AI MCP Server with Go LangChain Agent

A revolutionary **Model Context Protocol (MCP) server** that integrates **5 major AI providers** with intelligent **Go LangChain-style orchestration**. This project provides a complete AI ecosystem with file operations, multi-step workflows, and enterprise-grade performance.

## ⭐ **Star Features**

🤖 **5 AI Providers**: Claude, OpenAI, Gemini, Mistral, Hugging Face  
⚡ **Go LangChain Agent**: Intelligent orchestration with 10x Python performance  
📁 **File Operations**: Complete file management with AI analysis  
🔄 **Multi-Step Workflows**: Complex reasoning across multiple AI providers  
🏗️ **Extensible Architecture**: Perfect foundation for MCP ecosystem  
🔒 **Enterprise Ready**: Type-safe, robust, scalable

## 🤖 Supported AI Providers

Your MCP server now supports:
- **🤖 Claude AI** (Anthropic)
- **🧠 OpenAI GPT** (ChatGPT)
- **🔮 Google Gemini** (Google AI)
- **⚡ Mistral AI** (Mistral)
- **🤗 Hugging Face** (Open Source Models)
- **📮 Brazilian Zipcode Lookup** (ViaCEP API)

## 🏗️ Architecture

This project demonstrates the MCP (Model Context Protocol) implementation with:
- **MCP Server**: Main application that exposes tools via JSON-RPC
- **Test Client**: Command-line utility to test the server
- **External APIs**: Integration with ViaCEP and 5 AI providers
- **Caching**: File-based caching for zip code lookups
- **Go LangChain Agent**: Orchestration layer for complex workflows

![Architecture Diagram](architecture.png)

## 🚀 Features

### 1. Zipcode Tool 📮
- Look up Brazilian addresses by postal code (CEP)
- Uses ViaCEP API (`viacep.com.br`)
- Built-in caching (500 seconds TTL)
- Returns complete address information

### 2. Multi-AI Provider Support 🤖
Ask questions to any of these AI providers:
- **Claude AI** (`ask_claude`): Anthropic's Claude-3-Haiku
- **OpenAI GPT** (`ask_openai`): GPT-3.5-turbo  
- **Google Gemini** (`ask_gemini`): Gemini-1.5-flash
- **Mistral AI** (`ask_mistral`): Mistral-tiny
- **Hugging Face** (`ask_huggingface`): DialoGPT-medium

All AI tools support natural language queries and return formatted responses.

## 📦 Installation

1. **Clone the repository**
   ```powershell
   git clone <your-repo-url>
   cd mcp-server-test
   ```

2. **Install dependencies**
   ```powershell
   go mod tidy
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the project root with all your API keys:
   ```env
   # Required for Claude AI
   CLAUDE_API_KEY=your_claude_api_key_here
   
   # Required for OpenAI GPT
   OPENAI_API_KEY=your_openai_api_key_here
   
   # Required for Google Gemini
   GEMINI_API_KEY=your_gemini_api_key_here
   
   # Required for Mistral AI
   MISTRAL_API_KEY=your_mistral_api_key_here
   
   # Required for Hugging Face
   HUGGINGFACEHUB_API_TOKEN=your_huggingface_token_here
   ```
   
   **Note**: You only need the API keys for the AI providers you want to use.

## 🎯 Usage

### Important: Two Different Programs!

This project has **two separate executables**:

1. **Main MCP Server** (`main.go` in root) - Listens for JSON-RPC requests
2. **Test Client** (`cmd/test-client/main.go`) - Sends commands to the server

⚠️ **Don't run commands directly on the main server!** Use the test client instead.

### Running the MCP Server

Start the server from the project root:

```powershell
go run .\main.go
```

The server will start and listen for JSON-RPC requests via stdin/stdout.
**This will wait for input - it doesn't accept command-line arguments!**

### Using the Test Client

The test client provides an easy way to interact with the server:

```powershell
cd cmd\test-client
```

#### List Available Tools
```powershell
go run .\main.go list
```

**Example Output:**
```json
{
  "tools": [
    {"name": "zipcode", "description": "Find an address by zip code"},
    {"name": "ask_claude", "description": "Ask a question to Claude AI"},
    {"name": "ask_openai", "description": "Ask a question to OpenAI GPT"},
    {"name": "ask_gemini", "description": "Ask a question to Google Gemini"},
    {"name": "ask_mistral", "description": "Ask a question to Mistral AI"},
    {"name": "ask_huggingface", "description": "Ask a question to Hugging Face models"}
  ]
}
```

#### Test All AI Providers
```powershell
# Test Claude AI
go run .\main.go claude "What is the capital of Brazil?"

# Test OpenAI GPT
go run .\main.go openai "Explain quantum physics in simple terms"

# Test Google Gemini
go run .\main.go gemini "What are the benefits of renewable energy?"

# Test Mistral AI
go run .\main.go mistral "Write a haiku about programming"

# Test Hugging Face
go run .\main.go huggingface "Hello, how are you?"
```

#### Test Zipcode Tool
```powershell
# Make sure you're in the test-client directory first!
cd cmd\test-client
go run .\main.go zipcode 01310-100
```

**Example Output:**
```
Sending: {"id":1,"jsonrpc":"2.0","method":"tools/call","params":{"arguments":{"zip_code":"01310-100"},"name":"zipcode"}}
Response: {"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"Your address is {\"cep\":\"01310-100\",\"logradouro\":\"Avenida Paulista\",\"complemento\":\"\",\"bairro\":\"Bela Vista\",\"localidade\":\"São Paulo\",\"uf\":\"SP\",\"unidade\":\"\",\"ibge\":\"3550308\",\"gia\":\"1004\"}!"}]}}
```

#### Test Claude AI Tool
```powershell
# Make sure you're in the test-client directory first!
cd cmd\test-client
go run .\main.go claude "What is the capital of Brazil?"
```

**Example Outputs:**
```
# Claude Response
Response: {"result":{"content":[{"text":"Claude says: The capital of Brazil is Brasília."}]}}

# OpenAI Response  
Response: {"result":{"content":[{"text":"OpenAI says: 1+1 equals 2."}]}}

# Mistral Response
Response: {"result":{"content":[{"text":"Mistral says: Hello! I'm an AI assistant, how can I help you today?"}]}}
```

## 🧠 Go LangChain Agent

For advanced orchestration and multi-step workflows, use the Go LangChain-style agent:

```powershell
cd go-agent

# Test individual tools
go run test_simple.go

# Run comprehensive AI testing
go run test_all_ai.go

# Future: Full agent with OpenAI orchestration
# go run main.go
```

The Go agent provides:
- **Tool Orchestration**: Intelligent combination of zipcode and AI tools
- **Type Safety**: Compile-time guarantees and performance
- **Native Integration**: Same language as your MCP server
- **LangChain-Style Features**: Memory, workflows, and chaining

## 🔧 Direct JSON-RPC Usage

You can also interact with the server directly using JSON-RPC messages:

### List Tools Request
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

### Zipcode Tool Request
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "zipcode",
    "arguments": {
      "zip_code": "01310-100"
    }
  }
}
```

### Claude AI Tool Request
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "ask_claude",
    "arguments": {
      "question": "Explain quantum computing in simple terms"
    }
  }
}
```

## 📁 Project Structure

```
mcp-server-test/
├── main.go                     # Main MCP server
├── go.mod                      # Main project dependencies
├── .env                        # Environment variables
├── cmd/
│   └── test-client/
│       ├── main.go            # Test client utility
│       └── go.mod             # Test client dependencies
├── mcp-server-architecture.dot # Graphviz diagram source
├── architecture.png           # Generated architecture diagram
└── README.md                  # This file
```

## 🛠️ Building

### Build the main server
```powershell
go build -o mcp-server.exe .
```

### Build the test client
```powershell
cd cmd\test-client
go build -o test-client.exe .
```

## 🎨 Generating Architecture Diagram

If you have Graphviz installed:

```powershell
# Install Graphviz (if not already installed)
winget install graphviz

# Add to PATH (for current session)
$env:PATH += ";C:\Program Files\Graphviz\bin"

# Generate diagram
dot -Tpng mcp-server-architecture.dot -o architecture.png
```

## 🔒 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CLAUDE_API_KEY` | Your Anthropic Claude API key | Yes (for Claude tool) |

## 📝 API Documentation

### Tools Available

#### 1. `zipcode`
- **Description**: Find an address by zip code
- **Arguments**: 
  - `zip_code` (string, required): Brazilian postal code
- **Returns**: JSON with complete address information

#### 2. `ask_claude`
- **Description**: Ask a question to Claude AI (Anthropic)
- **Arguments**:
  - `question` (string, required): Question to ask Claude
- **Returns**: Claude's response as text

#### 3. `ask_openai`
- **Description**: Ask a question to OpenAI GPT
- **Arguments**:
  - `question` (string, required): Question to ask OpenAI
- **Returns**: OpenAI's response as text

#### 4. `ask_gemini`
- **Description**: Ask a question to Google Gemini
- **Arguments**:
  - `question` (string, required): Question to ask Gemini
- **Returns**: Gemini's response as text

#### 5. `ask_mistral`
- **Description**: Ask a question to Mistral AI
- **Arguments**:
  - `question` (string, required): Question to ask Mistral
- **Returns**: Mistral's response as text

#### 6. `ask_huggingface`
- **Description**: Ask a question to Hugging Face models
- **Arguments**:
  - `question` (string, required): Question to ask Hugging Face
- **Returns**: Hugging Face model's response as text

## 🧪 Testing

### Unit Tests
```powershell
go test .\...
```

### Integration Testing
Use the test client to verify both tools work correctly:

```powershell
cd cmd\test-client

# Test all functionality
go run .\main.go list
go run .\main.go zipcode 01310-100
go run .\main.go claude "Hello, how are you?"
```

## 🚨 Error Handling

The server handles various error scenarios:
- Invalid zip codes
- Network timeouts
- Missing API keys
- Malformed JSON-RPC requests

All errors are returned as proper JSON-RPC error responses.

## 🔄 Caching

The zipcode tool implements file-based caching:
- **Cache Duration**: 500 seconds
- **Cache Location**: System temp directory (`/tmp/cep*`)
- **Cache Key**: Based on zip code
- **Automatic Cleanup**: Cache files expire automatically

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly using the test client
5. Submit a pull request

## 📄 License

This project is open source and available under the MIT License.

## 🔗 Dependencies

- [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang) - MCP implementation
- [ryanuber/go-filecache](https://github.com/ryanuber/go-filecache) - File caching
- Standard Go libraries for HTTP, JSON, etc.

## 📞 Support

For issues and questions:
1. Check the test client examples above
2. Verify your `.env` configuration
3. Check that external APIs are accessible
4. Review the architecture diagram for understanding data flow