package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// Complete workflow examples and demos
func main() {
	fmt.Println("🚀 Go LangChain Agent - Complete Multi-AI Workflows")
	fmt.Println(strings.Repeat("=", 70))

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: OPENAI_API_KEY not found in environment")
		fmt.Println("Please set your OpenAI API key in the .env file")
		return
	}

	agent := NewGoLangChainAgent(apiKey, "../")
	ctx := context.Background()

	// Workflow Examples
	workflows := []struct {
		name        string
		description string
		query       string
	}{
		{
			name:        "AI Comparison Workflow",
			description: "Compare responses from multiple AI providers",
			query:       `Compare how different AI providers answer this question: "What is the future of artificial intelligence?" Use the ai_comparison tool to get perspectives from Claude, OpenAI, and Mistral.`,
		},
		{
			name:        "Location Intelligence Workflow",
			description: "Address lookup with AI context analysis",
			query:       `I want to learn about a location in Brazil. First, look up the address for zipcode 01310-100, then ask Claude to provide historical and cultural context about that specific street and neighborhood.`,
		},
		{
			name:        "Multi-Provider Research Workflow",
			description: "Research workflow using multiple AI providers",
			query:       `I'm researching machine learning trends. Ask OpenAI about current ML trends, ask Gemini about Google's perspective on AI development, and ask Mistral about open-source AI developments. Then summarize the key themes.`,
		},
		{
			name:        "Decision-Making Workflow",
			description: "Use different AI providers for different aspects",
			query:       `I'm planning to visit São Paulo. First get the address for 01310-100, then ask Claude about tourist attractions nearby, ask OpenAI about the best restaurants in that area, and ask Mistral about transportation options.`,
		},
		{
			name:        "Technical Analysis Workflow",
			description: "Multi-step technical analysis using AI",
			query:       `Explain quantum computing from different perspectives: Ask Claude for a scientific explanation, ask OpenAI for practical applications, and ask Gemini for Google's research in this field.`,
		},
	}

	// Run each workflow
	for i, workflow := range workflows {
		fmt.Printf("\n🔬 Workflow %d: %s\n", i+1, workflow.name)
		fmt.Printf("📝 Description: %s\n", workflow.description)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("🤖 Query: %s\n\n", workflow.query)

		startTime := time.Now()
		fmt.Println("🔄 Processing...")

		response, err := agent.Chat(ctx, workflow.query)

		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Printf("✅ Completed in %v\n", duration)
			fmt.Printf("📋 Result:\n%s\n", response)
		}

		fmt.Println("\n" + strings.Repeat("=", 70))

		// Reset memory between workflows for clean demos
		if i < len(workflows)-1 {
			agent.ResetMemory()
		}
	}

	// Interactive mode
	fmt.Println("\n🎮 Interactive Mode")
	fmt.Println("Try these example commands:")
	fmt.Println("• 'Compare AI responses about climate change'")
	fmt.Println("• 'Look up zipcode 01310-100 and tell me about the area'")
	fmt.Println("• 'Ask different AIs about the future of programming'")
	fmt.Println("• Type 'quit' to exit, 'reset' to clear memory")
	fmt.Println(strings.Repeat("-", 70))

	for {
		fmt.Print("\n💬 You: ")
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) == "quit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if strings.ToLower(input) == "reset" {
			agent.ResetMemory()
			fmt.Println("🧹 Memory cleared!")
			continue
		}

		if input == "" {
			continue
		}

		fmt.Println("\n🤖 Agent:")
		startTime := time.Now()
		response, err := agent.Chat(ctx, input)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Printf("✅ Response (took %v):\n%s\n", duration, response)
		}
	}
}
