package review

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/olbrichattila/qreview/internal/helpers"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
)

// newAws creates a new AWS q reviewer
func newAws(
	retr retriever.Retriever,
	prompt string,
	reporters []report.Reporter,
	commentOnPR bool,
) Reviewer {
	return &awsq{
		reporters:   reporters,
		retr:        retr,
		prompt:      prompt,
		commentOnPR: commentOnPR,
	}
}

type awsq struct {
	reporters   []report.Reporter
	retr        retriever.Retriever
	prompt      string
	commentOnPR bool
}

// AnalyzeCode returns the result of the analyzed code
func (a *awsq) AnalyzeCode(fileName string) error {
	content, err := a.retr.Get(fileName)
	if err != nil {
		return fmt.Errorf("Analyze code %w", err)
	}

	remappedContent, lineMap := helpers.SourceCodeLineRemap(content.FileContent)

	var stdout, stderr bytes.Buffer
	fmt.Println("executing q command")

	cmd := exec.Command("/usr/bin/q", "chat", "--no-interactive", a.prompt+remappedContent)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("output", out.String())
		fmt.Printf("stdout: %s\n", stdout.String())
		fmt.Printf("stderr: %s\n", stderr.String())
		return fmt.Errorf("cannot execute aws Q command, %w", err)
	}

	fmt.Println("executed q command")
	rawResponse := out.String()
	aiResponse := stripAnsiCodes(rawResponse)

	// todo do a PR analyzer and comment if line is provided in the response
	if a.commentOnPR {
		err = commentOnPRIfNecessary(fileName, aiResponse, content.DiffContent, lineMap)
		if err != nil {
			return err
		}
	}

	return generateReports(a.reporters, fileName, aiResponse)
}

// Summary implements Reviewer.
func (a *awsq) Summary() error {
	return summary(a.reporters)
}

// stripAnsiCodes removes ANSI color codes and formatting from the input string
func stripAnsiCodes(str string) string {
	// This regex matches ANSI escape codes for colors and formatting
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	result := ansiRegex.ReplaceAllString(str, "")

	// Remove other special formatting sequences
	result = regexp.MustCompile(`\[0m\[0m+`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[38;5;\d+m`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[39m`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[90m`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[92m`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[1m`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[22m`).ReplaceAllString(result, "")

	return result
}
