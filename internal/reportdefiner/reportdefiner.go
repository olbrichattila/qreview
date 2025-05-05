// Package reportdefiner creates the dependencies for the giver report from a struct
package reportdefiner

import (
	"fmt"
	"os"

	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
	"github.com/olbrichattila/qreview/internal/review"
	"gopkg.in/yaml.v2"
)

// ReviewerDefinitions contains multiple ReviewerDefinition
type ReviewerDefinitions []ReviewerDefinition

// ReviewerDefinition contains AI prompt, the retriever kind, which is file or diff and list of reporters, html, markdown...
type ReviewerDefinition struct {
	Prompt        string               `yaml:"prompt"`
	RetrieverKind retriever.Kind       `yaml:"retrieverKind"`
	CommentOnPr   bool                 `yaml:"commentOnPr"`
	Reporters     []ReporterDefinition `yaml:"reporters"`
}

// ReporterDefinition defines a reporter, for it's kind with folder and name, Folder may not required if reporter does not save
type ReporterDefinition struct {
	Kind   report.Kind `yaml:"kind"` // Replace with report.Kind if available
	Folder string
	Name   string `yaml:"name"`
}

var fileRetriever retriever.Retriever
var diffRetriever retriever.Retriever
var mixedRetriever retriever.Retriever

func Load(envManager env.EnvironmentManager, yamlFileName, reportFolder string) ([]review.Reviewer, error) {
	if !fileExists(yamlFileName) {
		return GetDefaultReviewers(envManager, reportFolder)
	}

	var defs ReviewerDefinitions

	data, err := os.ReadFile(yamlFileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &defs)

	for i := 0; i < len(defs); i++ {
		for x := 0; x < len(defs[i].Reporters); x++ {
			defs[i].Reporters[x].Folder = reportFolder
		}
	}

	return GetReviewers(envManager, defs)
}

func GetDefaultReviewers(envManager env.EnvironmentManager, reportFolder string) ([]review.Reviewer, error) {
	def := ReviewerDefinitions{
		// {
		// 	Prompt:        review.PromptReview,
		// 	RetrieverKind: FileRetriever,
		// 	Reporters: []ReporterDefinition{
		// 		{Kind: report.KindHTML, Folder: reportFolder, Name: "review"},
		// 		{Kind: report.KindMarkdown, Folder: reportFolder, Name: "review"},
		// 	},
		// },
		// {
		// 	Prompt:        review.PromptReviewChanges,
		// 	RetrieverKind: DiffRetriever,
		// 	CommentOnPr:   true,
		// 	Reporters: []ReporterDefinition{
		// 		{Kind: report.KindHTML, Folder: reportFolder, Name: "changes"},
		// 		{Kind: report.KindMarkdown, Folder: reportFolder, Name: "changes"},
		// 		{Kind: report.KindSave, Folder: reportFolder, Name: "changes"},
		// 	},
		// },
		{
			Prompt:        review.PromptReview,
			RetrieverKind: retriever.KindMixed,
			CommentOnPr:   true,
			Reporters: []ReporterDefinition{
				{Kind: report.KindHTML, Folder: reportFolder, Name: "changes"},
				{Kind: report.KindMarkdown, Folder: reportFolder, Name: "changes"},
				{Kind: report.KindSave, Folder: reportFolder, Name: "changes"},
			},
		},
		// {
		// 	Prompt:        review.PromptExplainCode,
		// 	RetrieverKind: FileRetriever,
		// 	Reporters: []ReporterDefinition{
		// 		{Kind: report.KindHTML, Folder: reportFolder, Name: "documentation"},
		// 		{Kind: report.KindMarkdown, Folder: reportFolder, Name: "documentation"},
		// 	},
		// },
	}

	return GetReviewers(envManager, def)
}

// GetReviewers returns with the pre-built reviewr list
func GetReviewers(envManager env.EnvironmentManager, reviewerDefinitions ReviewerDefinitions) ([]review.Reviewer, error) {
	err := initRetrievers(envManager)
	if err != nil {
		return nil, err
	}

	reviewers := []review.Reviewer{}
	for _, reviewerDefinition := range reviewerDefinitions {
		currentRetriever, err := getRetrievers(envManager, reviewerDefinition.RetrieverKind)
		if err != nil {
			return nil, err
		}

		currentReporters, err := getReporters(reviewerDefinition.Reporters)
		if err != nil {
			return nil, err
		}

		reviewers = append(
			reviewers,
			review.New(envManager, currentRetriever, currentReporters, reviewerDefinition.Prompt, reviewerDefinition.CommentOnPr),
		)
	}

	return reviewers, nil
}

func initRetrievers(envManager env.EnvironmentManager) error {
	var err error
	fileRetriever, err = retriever.NewFile(envManager)
	if err != nil {
		return err
	}

	diffRetriever, err = retriever.NewGitDiff(envManager)
	if err != nil {
		return err
	}

	mixedRetriever, err = retriever.NewMixed(fileRetriever, diffRetriever)
	if err != nil {
		return err
	}

	return nil
}

func getRetrievers(envManager env.EnvironmentManager, retrieverKind retriever.Kind) (retriever.Retriever, error) {
	switch retrieverKind {
	case retriever.KindFile:
		return fileRetriever, nil
	case retriever.KindDiff:
		return diffRetriever, nil
	case retriever.KindMixed:
		return mixedRetriever, nil
	default:
		return nil, fmt.Errorf("cannot determine retriever, %s", retrieverKind)
	}
}

func getReporters(reporterDefinitions []ReporterDefinition) ([]report.Reporter, error) {
	currentReporters := []report.Reporter{}
	for _, reporterDef := range reporterDefinitions {
		currentReporter, err := report.New(reporterDef.Kind, reporterDef.Folder, reporterDef.Name)
		if err != nil {
			return nil, err
		}

		currentReporters = append(currentReporters, currentReporter)
	}

	return currentReporters, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
