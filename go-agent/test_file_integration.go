package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func testFileIntegration() {
	fmt.Println("🔧 Testing File Operations Integration...")

	// Get the current working directory (should be go-agent)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Get the parent directory (project root)
	serverPath := filepath.Dir(currentDir)
	fmt.Printf("Server path: %s\n", serverPath)

	// Create a test agent
	fmt.Println("\n1. Creating Go LangChain agent with file operations...")
	agent := NewGoLangChainAgent("test-key", serverPath)

	// Test file list tool
	fmt.Println("\n2. Testing file list tool...")
	listTool := NewMCPFileListTool(serverPath)
	result, err := listTool.Execute(map[string]interface{}{
		"directory": ".",
	})
	if err != nil {
		fmt.Printf("❌ File list test failed: %v\n", err)
	} else {
		fmt.Printf("✅ File list test passed:\n%s\n", result)
	}

	// Test file write tool
	fmt.Println("\n3. Testing file write tool...")
	writeTool := NewMCPFileWriteTool(serverPath)
	testContent := "Hello from Go LangChain Agent!\nThis file was created by the integrated file operations system.\nTimestamp: " + fmt.Sprintf("%v", os.Getpid())
	result, err = writeTool.Execute(map[string]interface{}{
		"file_path": "test_integration.txt",
		"content":   testContent,
	})
	if err != nil {
		fmt.Printf("❌ File write test failed: %v\n", err)
	} else {
		fmt.Printf("✅ File write test passed:\n%s\n", result)
	}

	// Test file read tool
	fmt.Println("\n4. Testing file read tool...")
	readTool := NewMCPFileReadTool(serverPath)
	result, err = readTool.Execute(map[string]interface{}{
		"file_path": "test_integration.txt",
	})
	if err != nil {
		fmt.Printf("❌ File read test failed: %v\n", err)
	} else {
		fmt.Printf("✅ File read test passed:\n%s\n", result)
	}

	// Show all available tools
	fmt.Println("\n5. Available tools in agent:")
	for toolName := range agent.tools {
		fmt.Printf("   - %s\n", toolName)
	}

	fmt.Println("\n🎉 File operations integration test complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Test multi-server workflows")
	fmt.Println("2. Create AI + file operation combinations")
	fmt.Println("3. Build more complex orchestration examples")
}
