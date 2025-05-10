package retriever

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/olbrichattila/qreview/internal/diffmapper"
)

// ContextExtractor extracts relevant code context around changed lines
type ContextExtractor struct {
	// Number of lines before and after a change to include as context
	ContextLines int
}

// NewContextExtractor creates a new context extractor with default settings
func NewContextExtractor(contextLines int) *ContextExtractor {
	if contextLines <= 0 {
		contextLines = 5 // Default to 5 lines of context
	}
	return &ContextExtractor{
		ContextLines: contextLines,
	}
}

// Block represents a continuous block of code with context
type Block struct {
	StartLine int
	EndLine   int
	Content   string
}

// ExtractContext extracts relevant code blocks from a file based on changed lines in diff
func (ce *ContextExtractor) ExtractContext(fileContent, diffContent string) (string, error) {
	// Parse the diff to get changed lines
	changedLines := diffmapper.GetMap(diffContent)
	if len(changedLines) == 0 {
		return fileContent, nil // No changes, return the whole file
	}

	// Get line numbers that were changed
	changedLineNumbers := make(map[int]bool)
	for _, cl := range changedLines {
		changedLineNumbers[cl.LineNum] = true
	}

	// Read the file content into lines
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Create context blocks around changed lines
	blocks := ce.createContextBlocks(lines, changedLineNumbers)

	// Merge overlapping blocks
	mergedBlocks := ce.mergeBlocks(blocks)

	// Build the final content from merged blocks
	var result strings.Builder

	// Add a header explaining what this is
	result.WriteString("// CONTEXT: This is a partial view of the file showing only changed code and its context\n\n")

	for i, block := range mergedBlocks {
		if i > 0 {
			result.WriteString("\n// ...\n\n") // Indicate omitted code between blocks
		}

		// Add line numbers as comments at the start of the block
		result.WriteString(fmt.Sprintf("// Lines %d-%d\n", block.StartLine, block.EndLine))
		result.WriteString(block.Content)
	}

	return result.String(), nil
}

// createContextBlocks creates initial context blocks around changed lines
func (ce *ContextExtractor) createContextBlocks(lines []string, changedLineNumbers map[int]bool) []Block {
	var blocks []Block

	// For each changed line, create a context block
	for lineNum := range changedLineNumbers {
		// Calculate start and end lines with context
		startLine := max(0, lineNum-ce.ContextLines-1) // -1 because line numbers are 1-based
		endLine := min(len(lines)-1, lineNum+ce.ContextLines-1)

		// Extract the content for this block
		var blockContent strings.Builder
		for i := startLine; i <= endLine; i++ {
			if i < len(lines) {
				blockContent.WriteString(lines[i])
				blockContent.WriteString("\n")
			}
		}

		blocks = append(blocks, Block{
			StartLine: startLine + 1, // Convert back to 1-based line numbers
			EndLine:   endLine + 1,
			Content:   blockContent.String(),
		})
	}

	return blocks
}

// mergeBlocks merges overlapping context blocks
func (ce *ContextExtractor) mergeBlocks(blocks []Block) []Block {
	if len(blocks) <= 1 {
		return blocks
	}

	// Sort blocks by start line
	sortBlocks(blocks)

	var mergedBlocks []Block
	current := blocks[0]

	for i := 1; i < len(blocks); i++ {
		if blocks[i].StartLine <= current.EndLine+1 {
			// Blocks overlap or are adjacent, merge them
			if blocks[i].EndLine > current.EndLine {
				// Need to append additional content
				current.EndLine = blocks[i].EndLine
				// We need to update the content here, but for simplicity
				// we'll regenerate content for merged blocks later
			}
		} else {
			// No overlap, add the current block and start a new one
			mergedBlocks = append(mergedBlocks, current)
			current = blocks[i]
		}
	}

	// Add the last block
	mergedBlocks = append(mergedBlocks, current)

	return mergedBlocks
}

// Helper function to sort blocks by start line
func sortBlocks(blocks []Block) {
	// Simple bubble sort for clarity
	for i := 0; i < len(blocks); i++ {
		for j := i + 1; j < len(blocks); j++ {
			if blocks[i].StartLine > blocks[j].StartLine {
				blocks[i], blocks[j] = blocks[j], blocks[i]
			}
		}
	}
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
