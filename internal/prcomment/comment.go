// Package prcomment comments on gitHub, Gitlab.. what is implemented
package prcomment

import "github.com/olbrichattila/qreview/internal/env"

type Commenter interface {
	Comment(prURL, filePath string, comment string, lineNumber int) error
}

func New(env env.EnvironmentManager) (Commenter, error) {
	// Currently it supports only github, for others please add .env variable
	// and switch case between them here
	return newGitHub(env)
}
