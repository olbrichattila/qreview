package git

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	gitResponse := out.String()
	files := []string{}
	scanner := bufio.NewScanner(strings.NewReader(gitResponse))
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	return files, nil
}

func GetDiff(fileName string) (string, error) {
	cmd := exec.Command("git", "diff", fileName)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}
