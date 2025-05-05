package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/olbrichattila/qreview/cmd"
	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/parentsummary"
	"github.com/olbrichattila/qreview/internal/reportdefiner"
	"github.com/olbrichattila/qreview/internal/reviewparser"
)

const (
	reportTypeReview = "review"
	reportTypeDiff   = "difference"
	reportTypeDoc    = "documentation"
)

func main2() {
	data, err := os.ReadFile("report/2025/05/05/10_09/changes/test.php.md")
	if err != nil {
		fmt.Println(err)
		return
	}

	parsed := reviewparser.Parse(string(data))

	fmt.Println(parsed.Summary)
	for key, value := range parsed.Lines {
		fmt.Println(key, " --- ", value)
	}
}

func main() {
	envManager, err := env.NewDotEnv()
	if err != nil {
		printErrors(err)
		return
	}

	reportFolder := fmt.Sprintf("report/%s", time.Now().Format("2006/01/02/15_04"))
	reviewers, err := reportdefiner.Load(envManager, "definitions.yaml", reportFolder)
	if err != nil {
		printErrors(err)
		return
	}

	command, err := cmd.New(envManager, reviewers)
	if err != nil {
		printErrors(err)
		return
	}

	if err := command.Execute(); err != nil {
		printErrors(err)
	}

	err = parentsummary.Generate("report", reportFolder)
	if err != nil {
		printErrors(err)
		return
	}
}

func printErrors(err error) {
	for err != nil {
		fmt.Println(err)
		err = errors.Unwrap(err)
	}
}
