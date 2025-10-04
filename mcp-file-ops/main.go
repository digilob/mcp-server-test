package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// File operations MCP server
type ReadFileArguments struct {
	FilePath string `json:"file_path" jsonschema:"required,description=Path to the file to read"`
}

type WriteFileArguments struct {
	FilePath string `json:"file_path" jsonschema:"required,description=Path to the file to write"`
	Content  string `json:"content" jsonschema:"required,description=Content to write to the file"`
}

type ListFilesArguments struct {
	Directory string `json:"directory" jsonschema:"required,description=Directory path to list files"`
}

type SearchFilesArguments struct {
	Directory string `json:"directory" jsonschema:"required,description=Directory to search in"`
	Pattern   string `json:"pattern" jsonschema:"required,description=Search pattern or text to find"`
}

type FileInfoArguments struct {
	FilePath string `json:"file_path" jsonschema:"required,description=Path to get file information"`
}

func main() {
	fmt.Println("Starting MCP File Operations Server...")

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	// Register read file tool
	server.RegisterTool("read_file", "Read contents of a file", func(args ReadFileArguments) (*mcp_golang.ToolResponse, error) {
		content, err := os.ReadFile(args.FilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(content))), nil
	})

	// Register write file tool
	server.RegisterTool("write_file", "Write content to a file", func(args WriteFileArguments) (*mcp_golang.ToolResponse, error) {
		err := os.WriteFile(args.FilePath, []byte(args.Content), 0644)
		if err != nil {
			return nil, fmt.Errorf("error writing file: %v", err)
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Successfully wrote to %s", args.FilePath))), nil
	})

	// Register list files tool
	server.RegisterTool("list_files", "List files in a directory", func(args ListFilesArguments) (*mcp_golang.ToolResponse, error) {
		files, err := os.ReadDir(args.Directory)
		if err != nil {
			return nil, fmt.Errorf("error reading directory: %v", err)
		}

		var fileList []string
		for _, file := range files {
			if file.IsDir() {
				fileList = append(fileList, file.Name()+"/")
			} else {
				fileList = append(fileList, file.Name())
			}
		}

		result := strings.Join(fileList, "\n")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(result)), nil
	})

	// Register search files tool
	server.RegisterTool("search_files", "Search for text in files", func(args SearchFilesArguments) (*mcp_golang.ToolResponse, error) {
		var results []string

		err := filepath.Walk(args.Directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue on errors
			}

			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".txt") {
				content, err := os.ReadFile(path)
				if err == nil && strings.Contains(string(content), args.Pattern) {
					results = append(results, fmt.Sprintf("Found in: %s", path))
				}
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error searching files: %v", err)
		}

		if len(results) == 0 {
			results = append(results, "No matches found")
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(strings.Join(results, "\n"))), nil
	})

	// Register file info tool
	server.RegisterTool("file_info", "Get file information", func(args FileInfoArguments) (*mcp_golang.ToolResponse, error) {
		info, err := os.Stat(args.FilePath)
		if err != nil {
			return nil, fmt.Errorf("error getting file info: %v", err)
		}

		fileInfo := fmt.Sprintf(`File: %s
Size: %d bytes
Modified: %s
Is Directory: %t
Mode: %s`,
			info.Name(),
			info.Size(),
			info.ModTime().Format("2006-01-02 15:04:05"),
			info.IsDir(),
			info.Mode())

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fileInfo)), nil
	})

	server.Serve()
}
