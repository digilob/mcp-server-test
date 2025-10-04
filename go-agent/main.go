package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
)

// Tool represents a function that can be called by the agent
type Tool interface {
	Name() string
	Description() string
	Execute(args map[string]interface{}) (string, error)
}

// MCPZipcodeTool wraps the MCP server zipcode functionality
type MCPZipcodeTool struct {
	serverPath string
}

func NewMCPZipcodeTool(serverPath string) *MCPZipcodeTool {
	return &MCPZipcodeTool{serverPath: serverPath}
}

func (t *MCPZipcodeTool) Name() string {
	return "mcp_zipcode_lookup"
}

func (t *MCPZipcodeTool) Description() string {
	return `Look up Brazilian addresses by postal code (CEP). 
	Input should be a map with "zipcode" key containing a Brazilian postal code like '01310-100'.
	Returns complete address information including street, neighborhood, city, and state.`
}

func (t *MCPZipcodeTool) Execute(args map[string]interface{}) (string, error) {
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

	// Parse the response
	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)

		// Parse JSON response
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			content := result.String()
			if strings.Contains(content, "Your address is") {
				jsonStr := strings.TrimSuffix(strings.Split(content, "Your address is ")[1], "!")

				// Parse address JSON
				var address map[string]interface{}
				if err := json.Unmarshal([]byte(jsonStr), &address); err == nil {
					// Format nicely
					formatted := fmt.Sprintf(`📍 Address found for %s:
🏠 Street: %s
🏘️ Neighborhood: %s  
🏙️ City: %s
🗺️ State: %s
📮 ZIP: %s`,
						zipcode,
						getStringValue(address, "logradouro"),
						getStringValue(address, "bairro"),
						getStringValue(address, "localidade"),
						getStringValue(address, "uf"),
						getStringValue(address, "cep"))
					return formatted, nil
				}
			}
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// MCPClaudeTool wraps the MCP server Claude AI functionality
type MCPClaudeTool struct {
	serverPath string
}

func NewMCPClaudeTool(serverPath string) *MCPClaudeTool {
	return &MCPClaudeTool{serverPath: serverPath}
}

func (t *MCPClaudeTool) Name() string {
	return "mcp_claude_ai"
}

func (t *MCPClaudeTool) Description() string {
	return `Ask questions to Claude AI through the MCP server.
	Input should be a map with "question" key containing a clear question or request.
	Returns Claude's response to your question.`
}

func (t *MCPClaudeTool) Execute(args map[string]interface{}) (string, error) {
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

	// Parse the response
	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)

		// Parse JSON response
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			return result.String(), nil
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// MCPOpenAITool wraps the MCP server OpenAI functionality
type MCPOpenAITool struct {
	serverPath string
}

func NewMCPOpenAITool(serverPath string) *MCPOpenAITool {
	return &MCPOpenAITool{serverPath: serverPath}
}

func (t *MCPOpenAITool) Name() string {
	return "mcp_openai_ai"
}

func (t *MCPOpenAITool) Description() string {
	return `Ask questions to OpenAI GPT through the MCP server.`
}

func (t *MCPOpenAITool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "openai", question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing OpenAI query: %v", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			return result.String(), nil
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// MCPGeminiTool wraps the MCP server Gemini functionality
type MCPGeminiTool struct {
	serverPath string
}

func NewMCPGeminiTool(serverPath string) *MCPGeminiTool {
	return &MCPGeminiTool{serverPath: serverPath}
}

func (t *MCPGeminiTool) Name() string {
	return "mcp_gemini_ai"
}

func (t *MCPGeminiTool) Description() string {
	return `Ask questions to Google Gemini through the MCP server.`
}

func (t *MCPGeminiTool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "gemini", question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing Gemini query: %v", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			return result.String(), nil
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// MCPMistralTool wraps the MCP server Mistral functionality
type MCPMistralTool struct {
	serverPath string
}

func NewMCPMistralTool(serverPath string) *MCPMistralTool {
	return &MCPMistralTool{serverPath: serverPath}
}

func (t *MCPMistralTool) Name() string {
	return "mcp_mistral_ai"
}

func (t *MCPMistralTool) Description() string {
	return `Ask questions to Mistral AI through the MCP server.`
}

func (t *MCPMistralTool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "mistral", question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing Mistral query: %v", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			return result.String(), nil
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// MCPHuggingFaceTool wraps the MCP server Hugging Face functionality
type MCPHuggingFaceTool struct {
	serverPath string
}

func NewMCPHuggingFaceTool(serverPath string) *MCPHuggingFaceTool {
	return &MCPHuggingFaceTool{serverPath: serverPath}
}

func (t *MCPHuggingFaceTool) Name() string {
	return "mcp_huggingface_ai"
}

func (t *MCPHuggingFaceTool) Description() string {
	return `Ask questions to Hugging Face models through the MCP server.`
}

func (t *MCPHuggingFaceTool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	cmd := exec.Command("go", "run", "main.go", "huggingface", question)
	cmd.Dir = clientPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing Hugging Face query: %v", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Response:") {
		responsePart := strings.Split(outputStr, "Response:")[1]
		responsePart = strings.TrimSpace(responsePart)
		result := gjson.Get(responsePart, "result.content.0.text")
		if result.Exists() {
			return result.String(), nil
		}
	}

	return fmt.Sprintf("Raw output: %s", outputStr), nil
}

// AIComparisonTool compares responses from multiple AI providers
type AIComparisonTool struct {
	serverPath string
}

func NewAIComparisonTool(serverPath string) *AIComparisonTool {
	return &AIComparisonTool{serverPath: serverPath}
}

func (t *AIComparisonTool) Name() string {
	return "ai_comparison"
}

func (t *AIComparisonTool) Description() string {
	return `Compare responses from multiple AI providers.`
}

func (t *AIComparisonTool) Execute(args map[string]interface{}) (string, error) {
	question, ok := args["question"].(string)
	if !ok {
		return "", fmt.Errorf("question argument is required")
	}

	// Default providers
	providers := []string{"claude", "openai", "mistral"}

	// Check if custom providers specified
	if providerList, ok := args["providers"].([]interface{}); ok {
		providers = make([]string, len(providerList))
		for i, p := range providerList {
			providers[i] = p.(string)
		}
	}

	clientPath := filepath.Join(t.serverPath, "cmd", "test-client")
	results := make(map[string]string)

	// Query each provider
	for _, provider := range providers {
		cmd := exec.Command("go", "run", "main.go", provider, question)
		cmd.Dir = clientPath

		output, err := cmd.CombinedOutput()
		if err != nil {
			results[provider] = fmt.Sprintf("Error: %v", err)
		} else {
			outputStr := string(output)
			if strings.Contains(outputStr, "Response:") {
				responsePart := strings.Split(outputStr, "Response:")[1]
				responsePart = strings.TrimSpace(responsePart)
				result := gjson.Get(responsePart, "result.content.0.text")
				if result.Exists() {
					results[provider] = result.String()
				} else {
					results[provider] = "No response parsed"
				}
			} else {
				results[provider] = "No response found"
			}
		}
	}

	// Format comparison
	comparison := fmt.Sprintf("🤖 AI Comparison for: \"%s\"\n\n", question)
	for _, provider := range providers {
		comparison += fmt.Sprintf("📋 %s:\n%s\n\n", strings.ToUpper(provider), results[provider])
	}

	return comparison, nil
}

// ConversationMemory stores conversation history
type ConversationMemory struct {
	messages []openai.ChatCompletionMessage
}

func NewConversationMemory() *ConversationMemory {
	return &ConversationMemory{
		messages: make([]openai.ChatCompletionMessage, 0),
	}
}

func (m *ConversationMemory) AddMessage(role, content string) {
	m.messages = append(m.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (m *ConversationMemory) GetMessages() []openai.ChatCompletionMessage {
	return m.messages
}

func (m *ConversationMemory) Clear() {
	m.messages = make([]openai.ChatCompletionMessage, 0)
}

// GoLangChainAgent provides LangChain-like orchestration in Go
type GoLangChainAgent struct {
	client    *openai.Client
	tools     map[string]Tool
	memory    *ConversationMemory
	systemMsg string
}

func NewGoLangChainAgent(apiKey, serverPath string) *GoLangChainAgent {
	client := openai.NewClient(apiKey)

	agent := &GoLangChainAgent{
		client: client,
		tools:  make(map[string]Tool),
		memory: NewConversationMemory(),
		systemMsg: `You are a helpful AI orchestrator with access to multiple tools and AI providers.

Available tools:
- mcp_zipcode_lookup: Look up Brazilian addresses by postal code. Use {"zipcode": "01310-100"} format.
- mcp_claude_ai: Ask questions to Claude AI (Anthropic). Use {"question": "your question"} format.
- mcp_openai_ai: Ask questions to OpenAI GPT. Use {"question": "your question"} format.
- mcp_gemini_ai: Ask questions to Google Gemini. Use {"question": "your question"} format.
- mcp_mistral_ai: Ask questions to Mistral AI. Use {"question": "your question"} format.
- mcp_huggingface_ai: Ask questions to Hugging Face models. Use {"question": "your question"} format.
- ai_comparison: Compare responses from multiple AI providers. Use {"question": "your question", "providers": ["claude", "openai", "gemini"]} format.

Capabilities:
1. Single AI queries for specific use cases
2. Multi-AI comparisons for diverse perspectives  
3. Address lookup with AI context analysis
4. Multi-step workflows combining tools
5. Intelligent provider selection based on query type

Tool usage format:
[TOOL:tool_name]
{"argument": "value"}
[/TOOL]

Always provide clear, formatted responses and explain your reasoning when using multiple tools.`,
	}

	// Register all tools
	agent.RegisterTool(NewMCPZipcodeTool(serverPath))
	agent.RegisterTool(NewMCPClaudeTool(serverPath))
	agent.RegisterTool(NewMCPOpenAITool(serverPath))
	agent.RegisterTool(NewMCPGeminiTool(serverPath))
	agent.RegisterTool(NewMCPMistralTool(serverPath))
	agent.RegisterTool(NewMCPHuggingFaceTool(serverPath))
	agent.RegisterTool(NewAIComparisonTool(serverPath))

	return agent
}

func (a *GoLangChainAgent) RegisterTool(tool Tool) {
	a.tools[tool.Name()] = tool
}

func (a *GoLangChainAgent) Chat(ctx context.Context, message string) (string, error) {
	// Add user message to memory
	a.memory.AddMessage("user", message)

	// Prepare messages for OpenAI
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: a.systemMsg},
	}
	messages = append(messages, a.memory.GetMessages()...)

	// Get response from OpenAI
	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.1,
	})
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	assistantResponse := resp.Choices[0].Message.Content

	// Check if the response contains tool calls
	processedResponse, err := a.processToolCalls(assistantResponse)
	if err != nil {
		return "", fmt.Errorf("error processing tool calls: %v", err)
	}

	// Add assistant response to memory
	a.memory.AddMessage("assistant", processedResponse)

	return processedResponse, nil
}

func (a *GoLangChainAgent) processToolCalls(response string) (string, error) {
	// Simple tool call parser looking for [TOOL:name] ... [/TOOL] patterns
	for {
		startIdx := strings.Index(response, "[TOOL:")
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(response[startIdx:], "[/TOOL]")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx

		// Extract tool call
		_ = response[startIdx : endIdx+7] // toolCall (not used but shows extraction)
		toolNameStart := startIdx + 6
		toolNameEnd := strings.Index(response[toolNameStart:], "]")
		if toolNameEnd == -1 {
			break
		}
		toolNameEnd += toolNameStart

		toolName := response[toolNameStart:toolNameEnd]

		// Extract arguments
		argsStart := toolNameEnd + 2 // Skip "]\n"
		argsEnd := endIdx
		argsJSON := strings.TrimSpace(response[argsStart:argsEnd])

		// Parse arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return response, fmt.Errorf("error parsing tool arguments: %v", err)
		}

		// Execute tool
		tool, exists := a.tools[toolName]
		if !exists {
			return response, fmt.Errorf("unknown tool: %s", toolName)
		}

		result, err := tool.Execute(args)
		if err != nil {
			result = fmt.Sprintf("Error executing tool %s: %v", toolName, err)
		}

		// Replace tool call with result
		response = response[:startIdx] + result + response[endIdx+7:]
	}

	return response, nil
}

func (a *GoLangChainAgent) ResetMemory() {
	a.memory.Clear()
}

// Helper function
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "N/A"
}
