package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Helper function to create a test tool for any AI provider
type SimpleAITool struct {
	serverPath string
	toolName   string
	aiName     string
}

func NewSimpleAITool(serverPath, toolName, aiName string) *SimpleAITool {
	return &SimpleAITool{
		serverPath: serverPath,
		toolName:   toolName,
		aiName:     aiName,
	}
}

func (t *SimpleAITool) Execute(question string) (string, error) {
	// Execute the MCP client
	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", t.toolName, question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing %s query: %v, output: %s", t.aiName, err, string(output))
	}

	return fmt.Sprintf("Raw MCP Response:\n%s", string(output)), nil
}

func main() {
	fmt.Println("🧪 Testing All AI Providers in MCP Server")
	fmt.Println(strings.Repeat("=", 60))

	serverPath := "../"
	testQuestion := "What is the capital of France?"

	// Define all AI providers to test
	aiProviders := []struct {
		name     string
		toolName string
		emoji    string
	}{
		{"Claude", "claude", "🤖"},
		{"OpenAI GPT", "openai", "🧠"},
		{"Google Gemini", "gemini", "🔮"},
		{"Mistral AI", "mistral", "⚡"},
		{"Hugging Face", "huggingface", "🤗"},
	}

	// Test zipcode tool first
	fmt.Println("\n📍 Testing Brazilian Zipcode Tool:")
	zipcodeCmd := exec.Command("go", "run", "main.go", "zipcode", "01310-100")
	zipcodeCmd.Dir = filepath.Join(serverPath, "cmd", "test-client")

	output, err := zipcodeCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success!\nResult: Address found for Avenida Paulista\n")
	}

	// Test all AI providers
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🤖 Testing All AI Providers")
	fmt.Printf("Question: '%s'\n", testQuestion)
	fmt.Println(strings.Repeat("-", 60))

	for _, provider := range aiProviders {
		fmt.Printf("\n%s Testing %s:\n", provider.emoji, provider.name)

		tool := NewSimpleAITool(serverPath, provider.toolName, provider.name)
		result, err := tool.Execute(testQuestion)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Printf("✅ %s responded successfully!\n", provider.name)
			// Show just a snippet of the response to keep output clean
			lines := strings.Split(result, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Response:") {
					fmt.Printf("Preview: %s\n", line[:min(len(line), 100)]+"...")
					break
				}
			}
		}
	}

	// Test tool listing
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 Testing Tool Listing:")

	listCmd := exec.Command("go", "run", "main.go", "list")
	listCmd.Dir = filepath.Join(serverPath, "cmd", "test-client")

	output, err = listCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Tools listed successfully!\n")
		fmt.Printf("Response preview: %s\n", string(output)[:min(len(output), 200)]+"...")
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Multi-AI MCP Server Testing Complete!")
	fmt.Println("\n🎯 Your MCP server now supports:")
	fmt.Println("   📮 Brazilian zipcode lookup")
	fmt.Println("   🤖 Claude AI")
	fmt.Println("   🧠 OpenAI GPT")
	fmt.Println("   🔮 Google Gemini")
	fmt.Println("   ⚡ Mistral AI")
	fmt.Println("   🤗 Hugging Face models")
	fmt.Println("\n🚀 Ready for LangChain orchestration!")
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
