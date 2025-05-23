package review

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/olbrichattila/qreview/internal/helpers"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
)

/*
	This is the ollama client impelentation. On linux install:
	curl -fsSL https://ollama.com/install.sh | sh
	Pull a model like:
	ollama pull llama3

	Example usage: ollama run llama3 "What's the capital of France?"

*/

// NewOllama creates a new mock reviewer
func newOllama(
	retr retriever.Retriever,
	prompt string,
	reporters []report.Reporter,
	commentOnPR bool,
) Reviewer {
	return &ollama{
		reporters:   reporters,
		retr:        retr,
		prompt:      prompt,
		commentOnPR: commentOnPR,
	}
}

type ollama struct {
	reporters   []report.Reporter
	retr        retriever.Retriever
	prompt      string
	commentOnPR bool
}

// AnalyzeCode returns the result of the analyzed code
func (a *ollama) AnalyzeCode(fileName string) error {
	var err error
	content, err := a.retr.Get(fileName)
	if err != nil {
		return fmt.Errorf("Analyze code %w", err)
	}

	remappedContent, lineMap := helpers.SourceCodeLineRemap(content.FileContent)

	fmt.Println("executing ollama command")
	cmd := exec.Command("ollama", "run", "llama3", a.prompt+remappedContent)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot execute ollama command, %w", err)
	}
	fmt.Println("executed ollama command")

	aiResponse := out.String()

	// FAKE cached AI response
	// data, err := os.ReadFile("report/2025/05/05/10_09/changes/test.php.md")
	// if err != nil {
	// 	return err
	// }

	// aiResponse := string(data)

	// todo do a PR analizer and comment if line is provided in the response

	// TODO mapper
	if a.commentOnPR {
		err = commentOnPRIfNecessary(fileName, aiResponse, content.DiffContent, lineMap)
		if err != nil {
			return err
		}
	}
	return generateReports(a.reporters, fileName, aiResponse)

}

// Summary implements Reviewer.
func (a *ollama) Summary() error {
	return summary(a.reporters)
}
