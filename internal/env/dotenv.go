package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environment variable constants
const (
	EnvAIClient           = "AI_CLIENT"
	EnvFileExtensions     = "FILE_EXTENSIONS"
	EnvGithubToken        = "GITHUB_TOKEN"
	EnvAwsAccessKeyID     = "AWS_ACCESS_KEY_ID"
	EnvAwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	EnvAwsRegion          = "AWS_REGION"
	EnvQReviewAPIEndpoint = "QREVIEW_API_ENDPOINT"
	EnvContextLines       = "CONTEXT_LINES"
)

// NewDotEnv creates a new environment manager that loads from .env file
func NewDotEnv() (EnvironmentManager, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Could not load .env file %w", err)
		}
	}

	return &dotenv{}, nil
}

type dotenv struct{}

// Client returns the AI client to use
func (e *dotenv) Client() string {
	return os.Getenv(EnvAIClient)
}

// FileExtensions returns the file extensions to process
func (e *dotenv) FileExtensions() []string {
	return getEnvAsSlice(EnvFileExtensions, ",")
}

// GithubToken returns the GitHub token
func (e *dotenv) GithubToken() string {
	return os.Getenv(EnvGithubToken)
}

// AwsAccessKeyID returns the AWS access key ID
func (e *dotenv) AwsAccessKeyID() string {
	return os.Getenv(EnvAwsAccessKeyID)
}

// AwsSecretAccessKey returns the AWS secret access key
func (e *dotenv) AwsSecretAccessKey() string {
	return os.Getenv(EnvAwsSecretAccessKey)
}

// AwsRegion returns the AWS region
func (e *dotenv) AwsRegion() string {
	region := os.Getenv(EnvAwsRegion)
	if region == "" {
		return "us-east-1"
	}
	return region
}

// QReviewAPIEndpoint returns the QReview API endpoint
func (e *dotenv) QReviewAPIEndpoint() string {
	endpoint := os.Getenv(EnvQReviewAPIEndpoint)
	if endpoint == "" {
		return "http://localhost:3001"
	}
	return endpoint
}

// ContextLines returns the number of context lines to include around changed code
func (e *dotenv) ContextLines() int {
	return getEnvAsInt(EnvContextLines, 5)
}

// ShouldProcessFile checks if the file should be processed based on its extension
func (e *dotenv) ShouldProcessFile(fileName string) bool {
	extensions := e.FileExtensions()
	if len(extensions) == 0 {
		return true
	}

	for _, ext := range extensions {
		if ext != "" && strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}

	return false
}

// Helper functions
func getEnvAsSlice(key, sep string) []string {
	val := os.Getenv(key)
	if val == "" {
		return []string{}
	}
	return strings.Split(val, sep)
}

func getEnvAsInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}
