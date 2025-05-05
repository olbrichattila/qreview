package retriever

import "fmt"

func NewMixed(fileRetriever, diffRetriever Retriever) (Retriever, error) {
	if fileRetriever == nil || diffRetriever == nil {
		return nil, fmt.Errorf("one or both of the retrievers in NewMixed is nil")
	}

	return &mixed{
		fileRetriever: fileRetriever,
		diffRetriever: diffRetriever,
	}, nil
}

type mixed struct {
	fileRetriever Retriever
	diffRetriever Retriever
}

// Get implements Retriever.
func (m *mixed) Get(fileName string) (Result, error) {
	fileResult, err := m.fileRetriever.Get(fileName)
	if err != nil {
		return Result{}, err
	}

	diffResult, err := m.diffRetriever.Get(fileName)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Kind:        KindMixed,
		FileContent: fileResult.FileContent,
		DiffContent: diffResult.FileContent,
	}, nil
}
