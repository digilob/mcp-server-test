package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/ryanuber/go-filecache"
)

const cacheTime = 500

type MyFunctionsArguments struct {
	ZipCode string `json:"zip_code" jsonschema:"required,description=The zip code to be searched"`
}

type ClaudeArguments struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask Claude AI"`
}

type OpenAIArguments struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask OpenAI GPT"`
}

type GeminiArguments struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask Google Gemini"`
	Model    string `json:"model" jsonschema:"description=Gemini model to use: pro, flash, flash-lite (default: flash)"`
}

type MistralArguments struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask Mistral AI"`
}

type HuggingFaceArguments struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask Hugging Face models"`
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// OpenAI types
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Gemini types
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Mistral types
type MistralRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type MistralResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Hugging Face types
type HuggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type HuggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

// Cep is the brazilian postal code and address information
type Cep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

// loadEnv loads environment variables from .env file (fails silently if file doesn't exist)
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		// Silently ignore if .env file doesn't exist (normal in Docker)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove trailing semicolon if present
			value = strings.TrimSuffix(value, ";")
			os.Setenv(key, value)
		}
	}
}

func main() {
	// Try to load .env file for local development
	// This will fail silently in Docker, which is fine since
	// environment variables are passed via docker-compose
	loadEnv()

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport()) // Register zipcode tool
	err := server.RegisterTool("zipcode", "Find an address by his zip code", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		address, err := getCep(arguments.ZipCode)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Your address is %s!", address))), nil
	})
	if err != nil {
		panic(err)
	}

	// Register Claude AI tool
	err = server.RegisterTool("ask_claude", "Ask a question to Claude AI", func(arguments ClaudeArguments) (*mcp_golang.ToolResponse, error) {
		answer, err := askClaude(arguments.Question)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Claude says: %s", answer))), nil
	})
	if err != nil {
		panic(err)
	}

	// Register OpenAI GPT tool
	err = server.RegisterTool("ask_openai", "Ask a question to OpenAI GPT", func(arguments OpenAIArguments) (*mcp_golang.ToolResponse, error) {
		answer, err := askOpenAI(arguments.Question)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("OpenAI says: %s", answer))), nil
	})
	if err != nil {
		panic(err)
	}

	// Register Gemini tool
	err = server.RegisterTool("ask_gemini", "Ask a question to Google Gemini", func(arguments GeminiArguments) (*mcp_golang.ToolResponse, error) {
		answer, err := askGemini(arguments.Question)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Gemini says: %s", answer))), nil
	})
	if err != nil {
		panic(err)
	}

	// Register Mistral tool
	err = server.RegisterTool("ask_mistral", "Ask a question to Mistral AI", func(arguments MistralArguments) (*mcp_golang.ToolResponse, error) {
		answer, err := askMistral(arguments.Question)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Mistral says: %s", answer))), nil
	})
	if err != nil {
		panic(err)
	}

	// Register Hugging Face tool
	err = server.RegisterTool("ask_huggingface", "Ask a question to Hugging Face models", func(arguments HuggingFaceArguments) (*mcp_golang.ToolResponse, error) {
		answer, err := askHuggingFace(arguments.Question)
		if err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hugging Face says: %s", answer))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}

func getCep(id string) (string, error) {
	cached := getFromCache(id)
	if cached != "" {
		return cached, nil
	}
	req, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", id))
	if err != nil {
		return "", err
	}

	var c Cep
	err = json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		return "", err
	}
	res, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return saveOnCache(id, string(res)), nil
}

func askClaude(question string) (string, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return "", errors.New("CLAUDE_API_KEY not found in environment")
	}

	requestBody := ClaudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 1000,
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Claude API error %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	err = json.NewDecoder(resp.Body).Decode(&claudeResp)
	if err != nil {
		return "", err
	}

	if len(claudeResp.Content) > 0 {
		return claudeResp.Content[0].Text, nil
	}

	return "", errors.New("no response from Claude")
}

func askOpenAI(question string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY not found in environment")
	}

	requestBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens: 1000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&openaiResp)
	if err != nil {
		return "", err
	}

	if len(openaiResp.Choices) > 0 {
		return openaiResp.Choices[0].Message.Content, nil
	}

	return "", errors.New("no response from OpenAI")
}

func askGemini(question string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY not found in environment")
	}

	requestBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: question},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	err = json.NewDecoder(resp.Body).Decode(&geminiResp)
	if err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", errors.New("no response from Gemini")
}

func askMistral(question string) (string, error) {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return "", errors.New("MISTRAL_API_KEY not found in environment")
	}

	requestBody := MistralRequest{
		Model: "mistral-tiny",
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens: 1000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Mistral API error %d: %s", resp.StatusCode, string(body))
	}

	var mistralResp MistralResponse
	err = json.NewDecoder(resp.Body).Decode(&mistralResp)
	if err != nil {
		return "", err
	}

	if len(mistralResp.Choices) > 0 {
		return mistralResp.Choices[0].Message.Content, nil
	}

	return "", errors.New("no response from Mistral")
}

func askHuggingFace(question string) (string, error) {
	apiKey := os.Getenv("HUGGINGFACEHUB_API_TOKEN")
	if apiKey == "" {
		return "", errors.New("HUGGINGFACEHUB_API_TOKEN not found in environment")
	}

	requestBody := HuggingFaceRequest{
		Inputs: question,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Using a popular text generation model
	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Hugging Face API error %d: %s", resp.StatusCode, string(body))
	}

	var hfResp HuggingFaceResponse
	err = json.NewDecoder(resp.Body).Decode(&hfResp)
	if err != nil {
		return "", err
	}

	if len(hfResp) > 0 {
		return hfResp[0].GeneratedText, nil
	}

	return "", errors.New("no response from Hugging Face")
}

func getFromCache(id string) string {
	updater := func(path string) error {
		return errors.New("expired")
	}

	fc := filecache.New(getCacheFilename(id), cacheTime*time.Second, updater)

	fh, err := fc.Get()
	if err != nil {
		return ""
	}

	content, err := io.ReadAll(fh)
	if err != nil {
		return ""
	}

	return string(content)
}

func saveOnCache(id string, content string) string {
	updater := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write([]byte(content))
		return err
	}

	fc := filecache.New(getCacheFilename(id), cacheTime*time.Second, updater)

	_, err := fc.Get()
	if err != nil {
		return ""
	}

	return content
}

func getCacheFilename(id string) string {
	return os.TempDir() + "/cep" + strings.Replace(id, "-", "", -1)
}
