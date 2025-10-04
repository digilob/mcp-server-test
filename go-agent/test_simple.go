package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Simplified tool implementations for testing
type SimpleTool interface {
	Name() string
	Description() string
	Execute(args map[string]interface{}) (string, error)
}

type SimpleMCPZipcodeTool struct {
	serverPath string
}

func NewSimpleMCPZipcodeTool(serverPath string) *SimpleMCPZipcodeTool {
	return &SimpleMCPZipcodeTool{serverPath: serverPath}
}

func (t *SimpleMCPZipcodeTool) Name() string {
	return "mcp_zipcode_lookup"
}

func (t *SimpleMCPZipcodeTool) Description() string {
	return "Look up Brazilian addresses by postal code"
}

func (t *SimpleMCPZipcodeTool) Execute(args map[string]interface{}) (string, error) {
	zipcode, ok := args["zipcode"].(string)
	if !ok {
		return "", fmt.Errorf("zipcode argument is required")
	}

	// Clean the zipcode
	zipcode = strings.ReplaceAll(zipcode, "-", "")
	zipcode = strings.ReplaceAll(zipcode, " ", "")
	if len(zipcode) == 8 {
		zipcode = fmt.Sprintf("%s-%s", zipcode[:5], zipcode[5:])
	}

	// Execute the MCP client
	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "zipcode", zipcode)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing zipcode lookup: %v, output: %s", err, string(output))
	}

	return fmt.Sprintf("Raw MCP Response:\n%s", string(output)), nil
}

type SimpleMCPClaudeTool struct {
	serverPath string
}

func NewSimpleMCPClaudeTool(serverPath string) *SimpleMCPClaudeTool {
	return &SimpleMCPClaudeTool{serverPath: serverPath}
}

func (t *SimpleMCPClaudeTool) Name() string {
	return "mcp_claude_ai"
}

func (t *SimpleMCPClaudeTool) Description() string {
	return "Ask questions to Claude AI through MCP server"
}

func (t *SimpleMCPClaudeTool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	// Execute the MCP client
	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "claude", question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing Claude query: %v, output: %s", err, string(output))
	}

	return fmt.Sprintf("Raw MCP Response:\n%s", string(output)), nil
}

func main() {
	fmt.Println("🧪 Testing Go LangChain MCP Tools Integration")
	fmt.Println(strings.Repeat("=", 50))

	serverPath := "../"

	// Test Zipcode Tool
	fmt.Println("\n📍 Testing Zipcode Tool:")
	zipcodeTool := NewSimpleMCPZipcodeTool(serverPath)
	
	args := map[string]interface{}{
		"zipcode": "01310-100",
	}
	
	result, err := zipcodeTool.Execute(args)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success!\n%s\n", result)
	}

	// Test Claude Tool
	fmt.Println("\n🤖 Testing Claude Tool:")
	claudeTool := NewSimpleMCPClaudeTool(serverPath)
	
	args = map[string]interface{}{
		"question": "What is 2 + 2?",
	}
	
	result, err = claudeTool.Execute(args)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success!\n%s\n", result)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ Tool testing complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Set OPENAI_API_KEY in your .env file")
	fmt.Println("2. Run the full agent: go run main.go")
}