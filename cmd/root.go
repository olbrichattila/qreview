package cmd

import (
	"fmt"

	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/review"
	"github.com/olbrichattila/qreview/internal/source"
)

// New creates a new command line interpreter
func New(env env.EnvironmentManager, reviewers []review.Reviewer) (CommandInterpreter, error) {
	// validation

	if env == nil {
		return nil, fmt.Errorf("environment manager should not be nil")
	}

	for _, reviewer := range reviewers {
		if reviewer == nil {
			return nil, fmt.Errorf("reviewer list contains a nil value")
		}
	}

	newSource, err := source.New(env)
	if err != nil {
		return nil, err
	}

	return &comm{
		env:       env,
		reviewers: reviewers,
		source:    newSource,
	}, nil
}

type CommandInterpreter interface {
	Execute() error
}

type comm struct {
	env       env.EnvironmentManager
	reviewers []review.Reviewer
	source    source.Source
}

func (c *comm) Execute() error {
	files, err := c.source.GetFiles()
	if err != nil {
		return fmt.Errorf("failed to get files form git: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	for _, file := range files {
		if !c.hasExt(file) {
			continue
		}

		err := c.executeReview(file)
		if err != nil {
			return fmt.Errorf("execute, could not read file: %w", err)
		}
	}

	return c.generateReportSummary()
}

func (c *comm) hasExt(fileName string) bool {
	return c.env.ShouldProcessFile(fileName)
}

func (c *comm) executeReview(fileName string) error {
	for _, reviewer := range c.reviewers {
		fmt.Printf("Reviewing %s...\n", fileName)
		if err := reviewer.AnalyzeCode(fileName); err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}
	}

	return nil
}

func (c *comm) generateReportSummary() error {
	for _, reviewer := range c.reviewers {
		if err := reviewer.Summary(); err != nil {
			return err
		}
	}

	return nil
}
