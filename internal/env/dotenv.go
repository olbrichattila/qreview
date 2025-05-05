package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func NewDotEnv() (EnvironmentManager, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Could not load .env file %w", err)
		}
	}

	return &dotenv{}, nil
}

type dotenv struct {
}

// ShouldProcessFile implements EnvironmentManager.
func (d *dotenv) ShouldProcessFile(fileName string) bool {
	extensions := os.Getenv(fileExtensions)
	extensionParts := strings.Split(extensions, ",")
	for _, extensionPart := range extensionParts {
		ext := "." + strings.TrimSpace(strings.ToLower(extensionPart))
		fileExt := filepath.Ext(fileName)
		if strings.ToLower(fileExt) == strings.ToLower(ext) {
			return true
		}
	}

	return false
}

// GithubToken implements EnvironmentManager.
func (d *dotenv) GithubToken() string {
	return os.Getenv(githubToken)
}

// Client implements EnvironmentManager.
func (d *dotenv) Client() string {
	return strings.ToLower(os.Getenv(clientEnvName))
}
