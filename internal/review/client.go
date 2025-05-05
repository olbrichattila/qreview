// Package review analyses code
package review

import (
	cmdinterpreter "github.com/olbrichattila/qreview/internal/cmd-interpreter"
	"github.com/olbrichattila/qreview/internal/diffmapper"
	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/prcomment"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
	"github.com/olbrichattila/qreview/internal/reviewparser"
)

const (
	PromptReview         = "Review this code for bugs, performance, security issues, and suggest improvements. Use the following format for your comments: Line: <line number>: <review>\n\n"
	PromptReviewChanges  = "Review this git diff, do not compare changes, only review new lines. Use the following format for your comments: Line: <line number>: <review>\n\n"
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

func commentOnPRIfNecessary(filePath string, comment, diffContent string) error {

	if cmdinterpreter.HasFlag(cmdinterpreter.FlagGithubPR) &&
		cmdinterpreter.HasFlag(cmdinterpreter.FlagComment) &&
		prCommenterCache != nil {
		remap := false
		if diffContent != "" {
			diffmapper.GetMap(diffContent)
			remap = true
		}
		prURL, err := cmdinterpreter.Flag(cmdinterpreter.FlagGithubPR)
		if err != nil {
			return err
		}

		parsedReview := reviewparser.Parse(comment)

		err = prCommenterCache.Comment(prURL, filePath, parsedReview.Summary, 1)
		if err != nil {
			return err
		}
		for lineNr, lineComment := range parsedReview.Lines {
			mappedLineNr := lineNr
			if remap {
				mappedLineNr, err = diffmapper.GetClosestPrOffset(lineNr)
				if err != nil {
					return err
				}
			}

			err = prCommenterCache.Comment(prURL, filePath, lineComment, mappedLineNr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
