package retriever

import (
	"github.com/olbrichattila/qreview/internal/env"
	"github.com/olbrichattila/qreview/internal/source"
)

func NewFile(env env.EnvironmentManager) (Retriever, error) {
	source, err := source.New(env)
	if err != nil {
		return nil, err
	}

	return &fileR{
		source: source,
	}, nil
}

type fileR struct {
	source source.Source
}

// Get implements Retriever.
func (f *fileR) Get(fileName string) (Result, error) {
	content, err := f.source.GetFile(fileName)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Kind:        KindFile,
		FileContent: content,
	}, nil
}
