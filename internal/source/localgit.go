package source

import (
	"fmt"
	"os"

	"github.com/olbrichattila/qreview/internal/git"
)

func newLocalGit() (Source, error) {
	return &localGit{}, nil
}

type localGit struct {
}

// GetDiff implements Source.
func (g *localGit) GetDiff(fileName string) (string, error) {
	fmt.Println("Git diff called")
	return git.GetDiff(fileName)
}

// GetFile implements Source.
func (g *localGit) GetFile(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return string(content), nil
}

// GetFiles implements Source.
func (g *localGit) GetFiles() ([]string, error) {
	return git.GetStagedFiles()
}
