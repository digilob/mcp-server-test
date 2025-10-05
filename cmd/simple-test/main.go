package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("🚀 Testing MCP Multi-AI Server Locally")
	fmt.Println("======================================")

	// Test 1: Start server and test zipcode (non-AI test first)
	fmt.Println("\n1. Testing zipcode tool (non-AI)...")
	testZipcode()

	// Test 2: Test Claude AI
	fmt.Println("\n2. Testing Claude AI...")
	testClaude()

	// Test 3: Test OpenAI
	fmt.Println("\n3. Testing OpenAI...")
	testOpenAI()

	fmt.Println("\n✅ Local MCP Server tests completed!")
}

func testZipcode() {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "zipcode",
			"arguments": map[string]interface{}{
				"zip_code": "01310-100",
			},
		},
	}

	response := runMCPTest(req)
	fmt.Printf("Zipcode Response: %s\n", response)
}

func testClaude() {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_claude",
			"arguments": map[string]interface{}{
				"question": "Hello! Please respond with 'Claude is working' to test the system.",
			},
		},
	}

	response := runMCPTest(req)
	fmt.Printf("Claude Response: %s\n", response)
}

func testOpenAI() {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_openai",
			"arguments": map[string]interface{}{
				"question": "Hello! Please respond with 'OpenAI is working' to test the system.",
			},
		},
	}

	response := runMCPTest(req)
	fmt.Printf("OpenAI Response: %s\n", response)
}

func runMCPTest(request map[string]interface{}) string {
	// Create JSON request
	reqBytes, _ := json.Marshal(request)

	// Write to temporary file
	tmpFile := "temp_request.json"
	err := os.WriteFile(tmpFile, reqBytes, 0644)
	if err != nil {
		return fmt.Sprintf("Error writing request: %v", err)
	}
	defer os.Remove(tmpFile)

	// Run MCP server with the request
	cmd := exec.Command("cmd", "/c", "type "+tmpFile+" | go run main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running command: %v, Output: %s", err, string(output))
	}

	return string(output)
}
