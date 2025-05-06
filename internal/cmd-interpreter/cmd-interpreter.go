package cmdinterpreter

import (
	"fmt"
	"os"
	"strings"
)

const (
	FlagGithubPR = "githubpr" // GitHub PR have to be processed, follower by PR URL
	FlagComment  = "comment"  // Also comment on the PR, if not set then it will be a screen/report only review
)

func Arg(index int) (string, error) {
	if index < 0 {
		return "", fmt.Errorf("index cannot be negative")
	}

	i2 := 0
	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "-") {
			if i2 == index {
				return strings.TrimSpace(os.Args[i]), nil
			}
			i2++
		}
	}

	return "", fmt.Errorf("cannot find argument at %d", index)
}

func HasFlag(name string) bool {
	_, err := Flag(name)
	return err == nil
}

func Flag(name string) (string, error) {
	lowerCaseName := "-" + strings.ToLower(name)
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			key, value, _ := strings.Cut(os.Args[i], "=")
			if strings.ToLower(key) == lowerCaseName {
				return value, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find command line flag %s", name)
}
