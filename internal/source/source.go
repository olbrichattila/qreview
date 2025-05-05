// Package source gets the data to be reviewed, from source, like local GIT or github. etc.
package source

import (
	"fmt"

	cmdinterpreter "github.com/olbrichattila/qreview/internal/cmd-interpreter"
	"github.com/olbrichattila/qreview/internal/env"
)

// Source implement this interface for data sources
type Source interface {
	GetFiles() ([]string, error)
	GetFile(fileName string) (string, error)
	GetDiff(fileName string) (string, error)
}

func New(environment env.EnvironmentManager) (Source, error) {
	if cmdinterpreter.HasFlag(cmdinterpreter.FlagGithubPR) {
		return newFromCommandLine(environment)
	}

	return newLocalGit()
}

func newFromCommandLine(environment env.EnvironmentManager) (Source, error) {
	prURL, err := cmdinterpreter.Flag(cmdinterpreter.FlagGithubPR)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s flag from command line", cmdinterpreter.FlagGithubPR)
	}
	return newGitHub(environment, prURL)

}
