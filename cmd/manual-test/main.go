package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🧪 Testing MCP Server with Real API Keys")
	fmt.Println("========================================")

	// Start the server
	cmd := exec.Command("go", "run", "main.go")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Give server time to start
	time.Sleep(2 * time.Second)

	// Test 1: List tools
	fmt.Println("\n1. Testing tools/list...")
	listReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	reqBytes, _ := json.Marshal(listReq)
	fmt.Printf("Sending: %s\n", string(reqBytes))

	_, err = stdin.Write(append(reqBytes, '\n'))
	if err != nil {
		log.Printf("Error writing to stdin: %v", err)
	}

	// Read response
	response := make([]byte, 4096)
	n, err := stdout.Read(response)
	if err != nil {
		log.Printf("Error reading response: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(response[:n]))
	}

	// Test 2: Test Claude (if tools/list worked)
	time.Sleep(1 * time.Second)
	fmt.Println("\n2. Testing Claude AI...")
	claudeReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "ask_claude",
			"arguments": map[string]interface{}{
				"question": "Hello! This is a test of the MCP system. Please respond briefly.",
			},
		},
	}

	reqBytes2, _ := json.Marshal(claudeReq)
	fmt.Printf("Sending: %s\n", string(reqBytes2))

	_, err = stdin.Write(append(reqBytes2, '\n'))
	if err != nil {
		log.Printf("Error writing to stdin: %v", err)
	}

	// Read Claude response (may take longer)
	time.Sleep(3 * time.Second)
	response2 := make([]byte, 8192)
	n2, err := stdout.Read(response2)
	if err != nil {
		log.Printf("Error reading Claude response: %v", err)
	} else {
		fmt.Printf("Claude Response: %s\n", string(response2[:n2]))
	}

	// Clean up
	stdin.Close()
	cmd.Process.Kill()

	fmt.Println("\n✅ MCP Server test completed!")
}
