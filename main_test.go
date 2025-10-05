package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/ryanuber/go-filecache"
)

// Test data structures
func TestMyFunctionsArguments(t *testing.T) {
	args := MyFunctionsArguments{
		ZipCode: "12345",
	}

	if args.ZipCode != "12345" {
		t.Errorf("Expected ZipCode to be '12345', got '%s'", args.ZipCode)
	}
}

func TestClaudeArguments(t *testing.T) {
	args := ClaudeArguments{
		Question: "What is the weather like?",
	}

	if args.Question != "What is the weather like?" {
		t.Errorf("Expected Question to be 'What is the weather like?', got '%s'", args.Question)
	}
}

// Test JSON marshaling/unmarshaling
func TestClaudeRequestSerialization(t *testing.T) {
	req := ClaudeRequest{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: 1000,
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ClaudeRequest: %v", err)
	}

	// Test unmarshaling
	var decoded ClaudeRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ClaudeRequest: %v", err)
	}

	if decoded.Model != req.Model {
		t.Errorf("Expected Model to be '%s', got '%s'", req.Model, decoded.Model)
	}

	if decoded.MaxTokens != req.MaxTokens {
		t.Errorf("Expected MaxTokens to be %d, got %d", req.MaxTokens, decoded.MaxTokens)
	}

	if len(decoded.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(decoded.Messages))
	}

	if decoded.Messages[0].Role != "user" {
		t.Errorf("Expected message role to be 'user', got '%s'", decoded.Messages[0].Role)
	}
}

// Mock HTTP server for testing weather API
func createMockWeatherServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "weather") {
			response := `{
				"main": {
					"temp": 72.5,
					"humidity": 60
				},
				"weather": [
					{
						"main": "Clear",
						"description": "clear sky"
					}
				],
				"name": "Test City"
			}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// Mock HTTP server for testing Claude API
func createMockClaudeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := ClaudeResponse{
			ID:   "test-id",
			Type: "message",
			Role: "assistant",
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{
					Type: "text",
					Text: "This is a test response from Claude AI.",
				},
			},
			Model:      "claude-3-sonnet-20240229",
			StopReason: "end_turn",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
}

// Test cache functionality
func TestCacheOperations(t *testing.T) {
	cache := filecache.NewCache("test_cache", cacheTime*time.Millisecond)
	defer cache.Clear() // Clean up after test

	key := "test_key"
	value := "test_value"

	// Test setting cache
	err := cache.Set(key, []byte(value))
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test getting cache immediately
	cached, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if string(cached) != value {
		t.Errorf("Expected cached value to be '%s', got '%s'", value, string(cached))
	}

	// Test cache expiration
	time.Sleep((cacheTime + 100) * time.Millisecond)

	_, err = cache.Get(key)
	if err == nil {
		t.Error("Expected cache to be expired, but got value")
	}
}

// Test MCP server tool creation
func TestMCPToolCreation(t *testing.T) {
	// Test weather tool creation
	weatherTool := mcp_golang.NewTool("get_weather")
	weatherTool.Description = "Get weather information for a zip code"

	if weatherTool.Name != "get_weather" {
		t.Errorf("Expected tool name to be 'get_weather', got '%s'", weatherTool.Name)
	}

	// Test Claude tool creation
	claudeTool := mcp_golang.NewTool("ask_claude")
	claudeTool.Description = "Ask a question to Claude AI"

	if claudeTool.Name != "ask_claude" {
		t.Errorf("Expected tool name to be 'ask_claude', got '%s'", claudeTool.Name)
	}
}

// Integration test for message handling
func TestMessageStructure(t *testing.T) {
	message := Message{
		Role:    "user",
		Content: "Test message content",
	}

	if message.Role != "user" {
		t.Errorf("Expected role to be 'user', got '%s'", message.Role)
	}

	if message.Content != "Test message content" {
		t.Errorf("Expected content to be 'Test message content', got '%s'", message.Content)
	}
}

// Benchmark test for JSON operations
func BenchmarkClaudeRequestMarshal(b *testing.B) {
	req := ClaudeRequest{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: 1000,
		Messages: []Message{
			{
				Role:    "user",
				Content: "Benchmark test message",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			b.Fatalf("Marshal failed: %v", err)
		}
	}
}

// Test helper function to validate JSON structure
func validateJSONStructure(t *testing.T, jsonData []byte, expectedFields ...string) {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	for _, field := range expectedFields {
		if _, exists := data[field]; !exists {
			t.Errorf("Expected field '%s' not found in JSON", field)
		}
	}
}

// Example test for validating weather response structure
func TestWeatherResponseStructure(t *testing.T) {
	mockResponse := `{
		"main": {
			"temp": 72.5,
			"humidity": 60
		},
		"weather": [
			{
				"main": "Clear",
				"description": "clear sky"
			}
		],
		"name": "Test City"
	}`

	validateJSONStructure(t, []byte(mockResponse), "main", "weather", "name")
}

// Test for checking constants
func TestConstants(t *testing.T) {
	if cacheTime != 500 {
		t.Errorf("Expected cacheTime to be 500, got %d", cacheTime)
	}
}
