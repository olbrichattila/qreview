package review

import (
	"fmt"

	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
)

// newMock creates a new mock reviewer
func newMock(retr retriever.Retriever, prompt string, reporters []report.Reporter, commentOnPR bool) Reviewer {
	return &mock{
		reporters:   reporters,
		retr:        retr,
		prompt:      prompt,
		commentOnPR: commentOnPR,
	}
}

type mock struct {
	reporters   []report.Reporter
	retr        retriever.Retriever
	prompt      string
	result      string
	commentOnPR bool
}

// AnalyzeCode returns the result of the analyzed code
func (a *mock) AnalyzeCode(fileName string) error {
	content, err := a.retr.Get(fileName)
	if err != nil {
		return fmt.Errorf("-- MOCK error %s--\n Prompt: %s", err.Error(), a.prompt)
	}

	fakeContent := fmt.Sprintf("-- MOCK result --\n Content:\n%s\n\nPrompt: %s", content.FileContent, a.prompt)
	if a.commentOnPR {
		err = commentOnPRIfNecessary(fileName, "this is an automated test comment", "")
		if err != nil {
			return err
		}
	}

	return generateReports(a.reporters, fileName, fakeContent)
}

// Summary implements Reviewer.
func (a *mock) Summary() error {
	return summary(a.reporters)
}
