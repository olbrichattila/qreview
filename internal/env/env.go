package env

import (
	"os"
	"strconv"
	"strings"
)

// EnvironmentManager interface for environment variables
type EnvironmentManager interface {
	Client() string
	FileExtensions() []string
	GithubToken() string
	AwsAccessKeyID() string
	AwsSecretAccessKey() string
	AwsRegion() string
	QReviewAPIEndpoint() string
	ContextLines() int
}

// New creates a new environment manager
func New() EnvironmentManager {
	return &env{}
}

type env struct{}

// Client returns the AI client to use
func (e *env) Client() string {
	return os.Getenv("AI_CLIENT")
}

// FileExtensions returns the file extensions to process
func (e *env) FileExtensions() []string {
	extensions := os.Getenv("FILE_EXTENSIONS")
	if extensions == "" {
		return []string{}
	}

	return strings.Split(extensions, ",")
}

// GithubToken returns the GitHub token
func (e *env) GithubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

// AwsAccessKeyID returns the AWS access key ID
func (e *env) AwsAccessKeyID() string {
	return os.Getenv("AWS_ACCESS_KEY_ID")
}

// AwsSecretAccessKey returns the AWS secret access key
func (e *env) AwsSecretAccessKey() string {
	return os.Getenv("AWS_SECRET_ACCESS_KEY")
}

// AwsRegion returns the AWS region
func (e *env) AwsRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "us-east-1"
	}
	return region
}

// QReviewAPIEndpoint returns the QReview API endpoint
func (e *env) QReviewAPIEndpoint() string {
	endpoint := os.Getenv("QREVIEW_API_ENDPOINT")
	if endpoint == "" {
		return "http://localhost:3001"
	}
	return endpoint
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