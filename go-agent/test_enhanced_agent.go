package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Simple test for the enhanced agent
func testEnhancedAgent() {
	fmt.Println("🧪 Testing Enhanced Go LangChain Agent")
	fmt.Println(strings.Repeat("=", 50))

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: OPENAI_API_KEY not found")
		return
	}

	agent := NewGoLangChainAgent(apiKey, "../")
	ctx := context.Background()

	// Test 1: Simple AI query
	fmt.Println("\n🧠 Test 1: Simple AI Query")
	fmt.Println(strings.Repeat("-", 30))

	response, err := agent.Chat(ctx, `Ask Claude: "What is 2+2?"`)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response: %s\n", response)
	}

	// Test 2: Address lookup
	fmt.Println("\n📍 Test 2: Address Lookup")
	fmt.Println(strings.Repeat("-", 30))

	response, err = agent.Chat(ctx, `Look up the address for Brazilian zipcode 01310-100`)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response: %s\n", response)
	}

	// Test 3: AI Comparison
	fmt.Println("\n🤖 Test 3: AI Comparison")
	fmt.Println(strings.Repeat("-", 30))

	response, err = agent.Chat(ctx, `Use the ai_comparison tool to compare how Claude, OpenAI, and Mistral answer: "What is the capital of France?"`)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response: %s\n", response)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ Enhanced Agent Testing Complete!")
}

func main() {
	testEnhancedAgent()
}
