package source

import (
	"fmt"
	"strings"

	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/git"
	"github.com/olbrichattila/qreview/internal/pr"
)

var cachedDiffFiles []pr.FileDiff

func newGitHub(env env.EnvironmentManager, prURL string) (Source, error) {
	if prURL == "" {
		return nil, fmt.Errorf("the PR URL is missing")
	}

	if !git.IsValidGitHubPRURL(prURL) {
		return nil, fmt.Errorf("the PR URL is invalid %s, should look like https://github.com/user/repo/pull/123", prURL)
	}

	return &github{
		env:   env,
		pr:    pr.New(env),
		prURL: prURL,
	}, nil
}

type github struct {
	env   env.EnvironmentManager
	pr    pr.PullRequest
	prURL string
}

// GetDiff implements Source.
func (g *github) GetDiff(fileName string) (string, error) {
	var err error
	if cachedDiffFiles == nil {
		cachedDiffFiles, err = g.pr.GetPRFileDiffs(g.prURL)
		if err != nil {
			return "", err
		}
	}

	for _, f := range cachedDiffFiles {
		if f.Filename == fileName {
			normalizedCode := strings.ReplaceAll(f.Patch, "\r\n", "\n")
			return normalizedCode, nil
		}
	}

	return "", fmt.Errorf("diff %s file not found", fileName)
}

// GetFile implements Source.
func (g *github) GetFile(fileName string) (string, error) {
	result, err := g.pr.GetPRFileContent(g.prURL, fileName)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(result, "\r\n", "\n"), nil
}

// GetFiles implements Source.
func (g *github) GetFiles() ([]string, error) {
	return g.pr.GetPRFiles(g.prURL)
}
