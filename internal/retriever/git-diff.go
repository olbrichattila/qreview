package retriever

import (
	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/source"
)

func NewGitDiff(env env.EnvironmentManager) (Retriever, error) {
	source, err := source.New(env)
	if err != nil {
		return nil, err
	}
	return &diff{
		source: source,
	}, nil
}

type diff struct {
	source source.Source
}

// Get implements Retriever.
func (f *diff) Get(fileName string) (Result, error) {
	content, err := f.source.GetDiff(fileName)
	if err != nil {
		return Result{}, err
	}

	// It is deliberately assigned to file content, Diff content only apply for mixed
	// here we want to do the review on the diff directly
	return Result{
		Kind:        KindDiff,
		FileContent: content,
	}, nil
}
