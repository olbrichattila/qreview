package prcomment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/git"
)

func newGitHub(env env.EnvironmentManager) (Commenter, error) {
	if env == nil {
		return nil, fmt.Errorf("please provide github token in your environment: `GITHUB_TOKEN`")
	}

	return &github{
		env: env,
	}, nil
}

type github struct {
	env env.EnvironmentManager
}

// Comment implements Commenter.
func (g *github) Comment(prURL, filePath string, comment string, lineNumber int) error {
	githubToken := g.env.GithubToken()
	owner, repo, prNumber, err := git.GetPRInfo(prURL)
	if err != nil {
		return err
	}

	commitSHA, err := g.getPRHeadSHA(githubToken, owner, repo, prNumber)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, prNumber)

	body := map[string]interface{}{
		"body":      comment,
		"commit_id": commitSHA,
		"path":      filePath,
		"line":      lineNumber,
		"side":      "RIGHT",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("failed to post comment. Status: %s: File: %s\n", resp.Status, filePath)
}

// getPRHeadSHA fetches the head commit SHA of a GitHub PR
func (g *github) getPRHeadSHA(token, owner, repo string, prNumber int) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var result struct {
		Head struct {
			SHA string `json:"sha"`
		} `json:"head"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Head.SHA, nil
}
