package source

import (
	"fmt"
	"os"
	"strings"

	"github.com/olbrichattila/qreview/internal/git"
)

func newLocalGit() (Source, error) {
	return &localGit{}, nil
}

type localGit struct {
}

// GetDiff implements Source.
func (g *localGit) GetDiff(fileName string) (string, error) {
	result, err := git.GetDiff(fileName)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(result, "\r\n", "\n"), nil
}

// GetFile implements Source.
func (g *localGit) GetFile(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return strings.ReplaceAll(string(content), "\r\n", "\n"), nil
}

// GetFiles implements Source.
func (g *localGit) GetFiles() ([]string, error) {
	return git.GetStagedFiles()
}
