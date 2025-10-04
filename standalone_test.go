package main

import (
	"fmt"
	"os"
)

// Standalone test version that doesn't use stdio transport
// This allows you to test your functions directly without MCP protocol
func main() {
	fmt.Println("=== MCP Server Standalone Test ===")
	
	// Test 1: Test zipcode function
	fmt.Println("\n--- Testing zipcode function ---")
	zipcode := "01310-100" // Example Brazilian zipcode
	if len(os.Args) > 1 {
		zipcode = os.Args[1]
	}
	
	address, err := getCep(zipcode)
	if err != nil {
		fmt.Printf("Error getting address for %s: %v\n", zipcode, err)
	} else {
		fmt.Printf("Address for %s: %s\n", zipcode, address)
	}
	
	// Test 2: Test Claude function (only if API key is available)
	fmt.Println("\n--- Testing Claude AI function ---")
	
	// Check if Claude API key is set
	claudeKey := os.Getenv("CLAUDE_API_KEY")
	if claudeKey == "" {
		fmt.Println("CLAUDE_API_KEY not set, skipping Claude test")
	} else {
		question := "What is the capital of France?"
		if len(os.Args) > 2 {
			question = os.Args[2]
		}
		
		answer, err := askClaude(question)
		if err != nil {
			fmt.Printf("Error asking Claude '%s': %v\n", question, err)
		} else {
			fmt.Printf("Claude's answer to '%s': %s\n", question, answer)
		}
	}
	
	// Test 3: Test cache functions
	fmt.Println("\n--- Testing cache functions ---")
	testId := "test-cache-123"
	testContent := "This is test content for caching"
	
	// Save to cache
	result := saveOnCache(testId, testContent)
	if result != "" {
		fmt.Printf("Successfully saved to cache: %s\n", testId)
	} else {
		fmt.Printf("Failed to save to cache: %s\n", testId)
	}
	
	// Get from cache
	cached := getFromCache(testId)
	if cached != "" {
		fmt.Printf("Successfully retrieved from cache: %s\n", cached)
	} else {
		fmt.Printf("Failed to retrieve from cache or cache expired\n")
	}
	
	fmt.Println("\n=== Test Complete ===")
}