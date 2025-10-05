package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("🔑 API Key Verification Test")
	fmt.Println("============================")

	// Load .env file manually
	loadEnv()

	// Check each API key
	keys := map[string]string{
		"CLAUDE_API_KEY":           os.Getenv("CLAUDE_API_KEY"),
		"OPENAI_API_KEY":           os.Getenv("OPENAI_API_KEY"),
		"GEMINI_API_KEY":           os.Getenv("GEMINI_API_KEY"),
		"MISTRAL_API_KEY":          os.Getenv("MISTRAL_API_KEY"),
		"HUGGINGFACEHUB_API_TOKEN": os.Getenv("HUGGINGFACEHUB_API_TOKEN"),
	}

	for name, value := range keys {
		if value == "" {
			fmt.Printf("❌ %s: NOT SET\n", name)
		} else if len(value) < 10 {
			fmt.Printf("⚠️  %s: TOO SHORT (%d chars)\n", name, len(value))
		} else {
			// Show first 10 and last 4 characters for security
			masked := value[:10] + "..." + value[len(value)-4:]
			fmt.Printf("✅ %s: %s\n", name, masked)
		}
	}
}

// Copy the loadEnv function from main.go
func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	// Simple line-by-line parsing
	buffer := make([]byte, 1024)
	n, _ := file.Read(buffer)
	content := string(buffer[:n])

	lines := []string{}
	current := ""
	for _, char := range content {
		if char == '\n' || char == '\r' {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		equalPos := -1
		for i, char := range line {
			if char == '=' {
				equalPos = i
				break
			}
		}

		if equalPos > 0 {
			key := trimSpace(line[:equalPos])
			value := trimSpace(line[equalPos+1:])
			if len(value) > 0 && value[len(value)-1] == ';' {
				value = value[:len(value)-1]
			}
			os.Setenv(key, value)
		}
	}
	return nil
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
