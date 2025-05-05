// Package pr retrieves PR info from git
package pr

import (
	"github.com/olbrichattila/qreview/internal/env"
)

type PullRequest interface {
	GetPRFiles(prURL string) ([]string, error)
	GetPRFileContent(prURL, filePath string) (string, error)
	GetPRFileDiffs(prURL string) ([]FileDiff, error)
}

func New(env env.EnvironmentManager) PullRequest {
	// Currently it supports only github, for others please add .env variable
	// and switch case between them here
	return newGitHub(env)
}
