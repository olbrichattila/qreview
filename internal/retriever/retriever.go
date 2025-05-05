// Package retriever is an adapter which retrieves the text should be process by AI prompt reviewer
package retriever

// Retriever types
type Kind string

const (
	KindFile  Kind = "file"
	KindDiff  Kind = "diff"
	KindMixed Kind = "mixed"
)

// Result is the retriever result
type Result struct {
	Kind        Kind
	FileContent string
	DiffContent string
}

// Retriever implement this interface for each retriever
type Retriever interface {
	Get(fileName string) (Result, error)
}
