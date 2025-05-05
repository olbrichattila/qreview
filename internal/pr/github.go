package pr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/git"
)

func newGitHub(env env.EnvironmentManager) PullRequest {
	return &gitHubPr{env: env}
}

type gitHubPr struct {
	env env.EnvironmentManager
}

func (g *gitHubPr) GetPRFiles(prURL string) ([]string, error) {
	owner, repo, pullNumber, err := git.GetPRInfo(prURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, repo, pullNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.env.GithubToken())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Filename)
	}
	return fileNames, nil
}

func (g *gitHubPr) GetPRFileContent(prURL, filePath string) (string, error) {
	owner, repo, prNumber, err := git.GetPRInfo(prURL)
	if err != nil {
		return "", err
	}

	ref, err := g.getPRHeadSHA(g.env.GithubToken(), owner, repo, prNumber)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, filePath, ref)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Bearer "+g.env.GithubToken())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var contentResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(contentResp.Content)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

type FileDiff struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Patch    string `json:"patch,omitempty"`
}

func (g *gitHubPr) GetPRFileDiffs(prURL string) ([]FileDiff, error) {
	owner, repo, pullNumber, err := git.GetPRInfo(prURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, repo, pullNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.env.GithubToken())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var diffs []FileDiff
	if err := json.NewDecoder(resp.Body).Decode(&diffs); err != nil {
		return nil, err
	}
	return diffs, nil
}

// getPRHeadSHA fetches the head commit SHA of a GitHub PR
func (g *gitHubPr) getPRHeadSHA(token, owner, repo string, prNumber int) (string, error) {
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
