// Package helpers contains help functions
package helpers

import (
	"bufio"
	"strings"
)

// SourceCodeLineRemap strips blank lines out of the file, and map lines we keep back so AI will not hallucinate on line numbers
func SourceCodeLineRemap(content string) (string, map[int]int) {
	newContent := strings.Builder{}
	remap := make(map[int]int)
	originalLine := 0
	remappedLine := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		originalLine++

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		remap[remappedLine] = originalLine
		remappedLine++
		newContent.WriteString(line + "\n")
	}

	return newContent.String(), remap
}
