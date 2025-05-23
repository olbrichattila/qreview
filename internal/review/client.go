// Package review analyses code
package review

import (
	"fmt"

	cmdinterpreter "github.com/olbrichattila/qreview/internal/cmd-interpreter"
	"github.com/olbrichattila/qreview/internal/diffmapper"
	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/prcomment"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
	"github.com/olbrichattila/qreview/internal/reviewparser"
)

const (
	PromptReview         = "Review this code only for bugs and bad practices. Only comment issues and not good practices. Refer to the exact line number in the file. Use the following format for your comments: Line: <line number>: <review> Code starts here: "
	PromptExplainChanges = "Explain changes of the following diff:\n\n"
	PromptExplainCode    = "Explain what this code do:\n\n"

	clientQ       = "amazon_q"
	clientBedrock = "bedrock"
	clientOllama  = "ollama"
	clientMock    = "mock"
)

var prCommenterCache prcomment.Commenter

// Reviewer interface have to be implemented
type Reviewer interface {
	AnalyzeCode(filename string) error
	Summary() error
}

func New(
	env env.EnvironmentManager,
	retr retriever.Retriever,
	reporters []report.Reporter,
	prompt string,
	commentOnPR bool,
) Reviewer {
	envName := env.Client()
	// TODO error handling properly
	commenter, err := prcomment.New(env)
	if err == nil {
		prCommenterCache = commenter
	}

	switch envName {
	case clientQ:
		return newAws(retr, prompt, reporters, commentOnPR)
	case clientBedrock:
		return newBedrock(retr, prompt, reporters, commentOnPR)
	case clientOllama:
		return newOllama(retr, prompt, reporters, commentOnPR)
	case clientMock:
		return newMock(retr, prompt, reporters, commentOnPR)
	default:
		return newAws(retr, prompt, reporters, commentOnPR)
	}
}

func generateReports(reporters []report.Reporter, fileName, mdContent string) error {
	for _, reporter := range reporters {
		if reporter != nil {
			if err := reporter.Report(fileName, mdContent); err != nil {
				return err
			}
		}
	}

	return nil
}

func summary(reporters []report.Reporter) error {
	for _, reporter := range reporters {
		if reporter != nil {
			if err := reporter.Summary("index"); err != nil {
				return err
			}
		}
	}

	return nil
}

func commentOnPRIfNecessary(filePath string, comment, diffContent string, lineMap map[int]int) error {
	if cmdinterpreter.HasFlag(cmdinterpreter.FlagGithubPR) &&
		cmdinterpreter.HasFlag(cmdinterpreter.FlagComment) &&
		prCommenterCache != nil {
		remap := false
		var diffMap diffmapper.ChangedLines
		if diffContent != "" {
			diffMap = diffmapper.GetMap(diffContent)
			remap = true
		}
		prURL, err := cmdinterpreter.Flag(cmdinterpreter.FlagGithubPR)
		if err != nil {
			return err
		}

		parsedReview := reviewparser.Parse(comment)

		defaultLine := 1
		// This might rather be an error?
		if diffMap != nil {
			defaultLine = diffMap[0].LineNum
		}

		fmt.Printf("Commenting on PR File: %s, line %d\n", filePath, defaultLine)
		err = prCommenterCache.Comment(prURL, filePath, parsedReview.Summary, defaultLine)
		if err != nil {
			return err
		}
		for lineNr, lineComment := range parsedReview.Lines {
			mappedLineNr := lineMap[lineNr]
			if remap {
				mappedLineNr, err = diffmapper.GetClosestPrOffset(mappedLineNr)
				if err != nil {
					// skip for now
					continue
					// return err
				}
			}

			fmt.Printf("Commenting on PR File: %s, line %d\n", filePath, mappedLineNr)
			err = prCommenterCache.Comment(prURL, filePath, lineComment, mappedLineNr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
