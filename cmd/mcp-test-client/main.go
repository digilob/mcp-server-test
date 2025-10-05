package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// MCP Protocol requires initialization before tool calls
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  list - List available tools")
		fmt.Println("  zipcode <zipcode> - Test zipcode tool")
		fmt.Println("  claude <question> - Test Claude AI tool")
		fmt.Println("  openai <question> - Test OpenAI GPT tool")
		fmt.Println("  gemini <question> - Test Google Gemini tool")
		fmt.Println("  mistral <question> - Test Mistral AI tool")
		fmt.Println("  huggingface <question> - Test Hugging Face tool")
		return
	}

	command := os.Args[1]

	// Start the MCP server process
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "../../" // Set working directory to root
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	defer cmd.Process.Kill()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Step 1: Initialize the MCP connection
	if !initializeMCP(stdin, stdout) {
		fmt.Println("❌ Failed to initialize MCP connection")
		return
	}

	// Step 2: Execute the requested command
	switch command {
	case "list":
		listTools(stdin, stdout)
	case "zipcode":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a zipcode")
			return
		}
		testZipcode(stdin, stdout, os.Args[2])
	case "claude":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		testClaude(stdin, stdout, question)
	case "openai":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		testOpenAI(stdin, stdout, question)
	case "gemini":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		testGemini(stdin, stdout, question)
	case "mistral":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		testMistral(stdin, stdout, question)
	case "huggingface":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		testHuggingFace(stdin, stdout, question)
	default:
		fmt.Println("Unknown command:", command)
	}
}

// Initialize MCP connection with proper handshake
func initializeMCP(stdin io.WriteCloser, stdout io.ReadCloser) bool {
	// Send initialize request
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "1.0",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"clientInfo": map[string]interface{}{
				"name":    "mcp-test-client",
				"version": "1.0.0",
			},
		},
	}

	fmt.Println("🔄 Initializing MCP connection...")
	response := sendRequestWithResponse(stdin, stdout, initRequest)

	if response != "" && strings.Contains(response, "\"result\"") {
		fmt.Println("✅ MCP connection initialized successfully")

		// Send notifications/initialized (required by MCP protocol)
		initNotification := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "notifications/initialized",
			"params":  map[string]interface{}{},
		}

		requestBytes, _ := json.Marshal(initNotification)
		stdin.Write(requestBytes)
		stdin.Write([]byte("\n"))

		fmt.Println("📤 Sent initialization notification")
		time.Sleep(500 * time.Millisecond) // Give server time to process

		return true
	} else {
		fmt.Printf("❌ Initialization failed: %s\n", response)
		return false
	}
}

func listTools(stdin io.WriteCloser, stdout io.ReadCloser) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	fmt.Println("📋 Listing available tools...")
	sendRequestWithResponse(stdin, stdout, request)
}

func testZipcode(stdin io.WriteCloser, stdout io.ReadCloser, zipcode string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "zipcode",
			"arguments": map[string]interface{}{
				"zip_code": zipcode,
			},
		},
	}

	fmt.Printf("📮 Testing zipcode lookup for: %s\n", zipcode)
	sendRequestWithResponse(stdin, stdout, request)
}

func testClaude(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_claude",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	fmt.Printf("🤖 Testing Claude AI with question: %s\n", question)
	sendRequestWithResponse(stdin, stdout, request)
}

func testOpenAI(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      4,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_openai",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	fmt.Printf("🧠 Testing OpenAI with question: %s\n", question)
	sendRequestWithResponse(stdin, stdout, request)
}

func testGemini(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      5,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_gemini",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	fmt.Printf("🔮 Testing Gemini with question: %s\n", question)
	sendRequestWithResponse(stdin, stdout, request)
}

func testMistral(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      6,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_mistral",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	fmt.Printf("⚡ Testing Mistral with question: %s\n", question)
	sendRequestWithResponse(stdin, stdout, request)
}

func testHuggingFace(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      7,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_huggingface",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	fmt.Printf("🤗 Testing Hugging Face with question: %s\n", question)
	sendRequestWithResponse(stdin, stdout, request)
}

func sendRequestWithResponse(stdin io.WriteCloser, stdout io.ReadCloser, request map[string]interface{}) string {
	// Send request
	requestBytes, _ := json.Marshal(request)
	fmt.Printf("📤 Sending: %s\n", string(requestBytes))

	_, err := stdin.Write(requestBytes)
	if err != nil {
		fmt.Printf("❌ Error writing request: %v\n", err)
		return ""
	}
	_, err = stdin.Write([]byte("\n"))
	if err != nil {
		fmt.Printf("❌ Error writing newline: %v\n", err)
		return ""
	}

	// Read response with timeout
	scanner := bufio.NewScanner(stdout)

	// Wait for response with timeout
	done := make(chan string, 1)
	go func() {
		for scanner.Scan() {
			response := scanner.Text()
			if strings.TrimSpace(response) != "" {
				done <- response
				return
			}
		}
		done <- ""
	}()

	select {
	case response := <-done:
		if response != "" {
			fmt.Printf("📥 Response: %s\n", response)
			return response
		} else {
			fmt.Println("❌ No response received")
			return ""
		}
	case <-time.After(10 * time.Second):
		fmt.Println("❌ Timeout waiting for response")
		return ""
	}
}
