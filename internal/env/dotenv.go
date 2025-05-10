package env

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NewDotEnv creates a new environment manager that loads from .env file
func NewDotEnv() (EnvironmentManager, error) {
	err := loadEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	return &env{}, nil
}

func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	return scanner.Err()
}

// ShouldProcessFile checks if the file should be processed based on its extension
func (e *env) ShouldProcessFile(fileName string) bool {
	extensions := e.FileExtensions()
	if len(extensions) == 0 {
		return true
	}

	for _, ext := range extensions {
		if strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}

	return false
}

// ContextLines returns the number of context lines to include around changed code
func (e *env) ContextLines() int {
	contextLines := os.Getenv("CONTEXT_LINES")
	if contextLines == "" {
		return 5 // Default to 5 lines of context
	}
	
	lines, err := strconv.Atoi(contextLines)
	if err != nil {
		return 5 // Default to 5 lines if parsing fails
	}
	
	return lines
}