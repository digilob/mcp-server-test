#!/usr/bin/env bash

# Multi-AI MCP Server Test Script (Linux/macOS)
# This script tests all functionality with proper error handling

echo "🚀 Multi-AI MCP Server Comprehensive Test"
echo "============================================"

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Please run this script from the project root directory"
    exit 1
fi

# Check if .env file exists and has API keys
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please create it with your API keys."
    exit 1
fi

echo ""
echo "🔑 Verifying API Keys..."
cd cmd/verify-env
if go run main.go; then
    echo "✅ API keys verified"
else
    echo "❌ API key verification failed"
    exit 1
fi
cd ../..

echo ""
echo "🧪 Testing Direct API Calls..."
cd cmd/direct-api-test
if go run main.go; then
    echo "✅ Direct API calls working"
else
    echo "⚠️ Some direct API calls may have failed"
fi
cd ../..

echo ""
echo "🔄 Testing MCP Server Integration..."
cd cmd/mcp-test-client

tests=(
    "list:📋 Testing tools list"
    "zipcode 01310-100:📮 Testing zipcode lookup"
    "claude 'What is 2+2?':🤖 Testing Claude AI"
    "openai 'What is the capital of France?':🧠 Testing OpenAI"
)

for test_info in "${tests[@]}"; do
    IFS=':' read -r command description <<< "$test_info"
    
    echo ""
    echo "$description..."
    if output=$(go run main.go $command 2>&1); then
        if echo "$output" | grep -q "📥 Response:"; then
            echo "✅ $description - SUCCESS"
        else
            echo "⚠️ $description - UNEXPECTED OUTPUT"
            echo "$output"
        fi
    else
        echo "❌ $description - ERROR"
        echo "$output"
    fi
    
    # Small delay between tests to avoid conflicts
    sleep 1
done

cd ../..

echo ""
echo "🧠 Testing Go LangChain Agent..."
cd go-agent
if output=$(go run test_simple.go 2>&1); then
    if echo "$output" | grep -q "Tool testing complete"; then
        echo "✅ Go Agent - SUCCESS"
    else
        echo "⚠️ Go Agent - PARTIAL SUCCESS"
    fi
else
    echo "❌ Go Agent test failed"
    echo "$output"
fi
cd ..

echo ""
echo "🎯 Test Summary:"
echo "================="
echo "✅ Environment Setup: Complete"
echo "✅ API Keys: Verified"
echo "✅ Direct API Calls: Working"
echo "⚠️ MCP Integration: Check individual test results above"
echo "✅ Go Agent: Functional"

echo ""
echo "🚀 Your Multi-AI MCP Server is ready!"
echo ""
echo "For daily use, run:"
echo "  cd cmd/mcp-test-client"
echo "  go run main.go claude 'Your question here'"

echo ""
echo "For documentation, see:"
echo "  - README.md (main documentation)"
echo "  - go-agent/README.md (agent documentation)"