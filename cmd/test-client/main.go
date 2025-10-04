package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Simple test client to interact with the MCP server
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

	// Start the MCP server process (from the root directory)
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

	// Create JSON-RPC messages based on command
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

	cmd.Process.Kill()
}

func listTools(stdin io.WriteCloser, stdout io.ReadCloser) {
	// Send list tools request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	sendRequest(stdin, stdout, request)
}

func testZipcode(stdin io.WriteCloser, stdout io.ReadCloser, zipcode string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "zipcode",
			"arguments": map[string]interface{}{
				"zip_code": zipcode,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func testClaude(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_claude",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func testOpenAI(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_openai",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func testGemini(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_gemini",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func testMistral(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_mistral",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func testHuggingFace(stdin io.WriteCloser, stdout io.ReadCloser, question string) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_huggingface",
			"arguments": map[string]interface{}{
				"question": question,
			},
		},
	}

	sendRequest(stdin, stdout, request)
}

func sendRequest(stdin io.WriteCloser, stdout io.ReadCloser, request map[string]interface{}) {
	// Send request
	requestBytes, _ := json.Marshal(request)
	fmt.Printf("Sending: %s\n", string(requestBytes))

	stdin.Write(requestBytes)
	stdin.Write([]byte("\n"))
	stdin.Close()

	// Read response
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		response := scanner.Text()
		if response != "" {
			fmt.Printf("Response: %s\n", response)
			break
		}
	}
}
