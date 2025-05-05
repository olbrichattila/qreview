package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// GetPRInfo phrases GitHub PR URL: returns owner, repo, prNumber, error
func GetPRInfo(prURL string) (string, string, int, error) {
	if !IsValidGitHubPRURL(prURL) {
		return "", "", 0, fmt.Errorf("the PR url is invalid: %s", prURL)
	}

	u, err := url.Parse(prURL)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "pull" {
		return "", "", 0, fmt.Errorf("not a valid PR URL %s", prURL)
	}

	owner := parts[0]
	repo := parts[1]
	prNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("not a valid PR URL, PR number is not a number %s", prURL)
	}

	return owner, repo, prNumber, nil
}

// IsValidGitHubPRURL checks whether the given URL is a valid GitHub Pull Request URL
func IsValidGitHubPRURL(prURL string) bool {
	re := regexp.MustCompile(`^https?://github\.com/[^/]+/[^/]+/pull/\d+$`)
	return re.MatchString(prURL)
}
