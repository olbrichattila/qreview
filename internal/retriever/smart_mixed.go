package retriever

import (
	"fmt"
)

// NewSmartMixed creates a retriever that intelligently combines file and diff information
// to provide only the relevant parts of a file for review
func NewSmartMixed(fileRetriever, diffRetriever Retriever, contextLines int) (Retriever, error) {
	if fileRetriever == nil || diffRetriever == nil {
		return nil, fmt.Errorf("one or both of the retrievers in NewSmartMixed is nil")
	}

	return &smartMixed{
		fileRetriever:    fileRetriever,
		diffRetriever:    diffRetriever,
		contextExtractor: NewContextExtractor(contextLines),
	}, nil
}

type smartMixed struct {
	fileRetriever    Retriever
	diffRetriever    Retriever
	contextExtractor *ContextExtractor
}

// Get implements Retriever.
func (m *smartMixed) Get(fileName string) (Result, error) {
	// Get the full file content
	fileResult, err := m.fileRetriever.Get(fileName)
	if err != nil {
		return Result{}, err
	}

	// Get the diff content
	diffResult, err := m.diffRetriever.Get(fileName)
	if err != nil {
		return Result{}, err
	}

	// Extract only the relevant parts of the file based on the diff
	contextContent, err := m.contextExtractor.ExtractContext(fileResult.FileContent, diffResult.FileContent)
	if err != nil {
		// Fall back to full file content if extraction fails
		return Result{
			Kind:        KindMixed,
			FileContent: fileResult.FileContent,
			DiffContent: diffResult.FileContent,
		}, nil
	}

	return Result{
		Kind:        KindMixed,
		FileContent: contextContent, // Only the relevant parts with context
		DiffContent: diffResult.FileContent,
	}, nil
}