package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("🧪 Direct API Test")
	fmt.Println("==================")

	// Load environment variables
	loadEnv()

	// Test Claude API directly
	fmt.Println("\n1. Testing Claude API directly...")
	testClaudeDirect()

	// Test OpenAI API directly
	fmt.Println("\n2. Testing OpenAI API directly...")
	testOpenAIDirect()
}

func testClaudeDirect() {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ CLAUDE_API_KEY not found")
		return
	}

	requestBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307",
		"max_tokens": 100,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "Hello! Please respond with 'Claude API is working' to test the connection.",
			},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", resp.Status)
	if resp.StatusCode == 200 {
		fmt.Printf("✅ Claude Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Claude Error: %s\n", string(body))
	}
}

func testOpenAIDirect() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ OPENAI_API_KEY not found")
		return
	}

	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "Hello! Please respond with 'OpenAI API is working' to test the connection.",
			},
		},
		"max_tokens": 100,
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", resp.Status)
	if resp.StatusCode == 200 {
		fmt.Printf("✅ OpenAI Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ OpenAI Error: %s\n", string(body))
	}
}

// Simplified loadEnv function
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) == 2 {
			key := string(bytes.TrimSpace(parts[0]))
			value := string(bytes.TrimSpace(parts[1]))
			os.Setenv(key, value)
		}
	}
}
