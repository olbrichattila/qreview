// Package diffmapper maps git diff with original file line number
package diffmapper

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type ChangedLines []ChangedLine

type ChangedLine struct {
	LineNum int
	Content string
}

var latestChanges ChangedLines

func GetMap(diff string) ChangedLines {
	hunkStarted := false
	var changes []ChangedLine
	scanner := bufio.NewScanner(strings.NewReader(diff))

	var newLineNum int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "@@") {
			fmt.Println("hunk", line)
			hunkStarted = true
			// Parse hunk header: @@ -a,b +c,d @@
			parts := strings.Split(line, " ")
			newHunk := parts[2] // "+c,d"
			hunkStart := strings.Split(strings.TrimPrefix(newHunk, "+"), ",")
			newLineNum, _ = strconv.Atoi(hunkStart[0])
			continue
		}

		if !hunkStarted {
			continue
		}

		switch {
		case strings.HasPrefix(line, "+"):
			changes = append(changes, ChangedLine{
				LineNum: newLineNum,
				Content: line[1:], // strip "+"
			})
			newLineNum++

		case strings.HasPrefix(line, "-"):
			// removed line, don't increment newLineNum
		case strings.HasPrefix(line, " "):
			newLineNum++
		}
	}

	latestChanges = changes
	return changes
}

func GetClosestPrOffset(prLineNr int) (int, error) {
	if latestChanges == nil {
		return 0, fmt.Errorf("you must run GetMap before getting the closest pr offset")
	}

	for i := len(latestChanges) - 1; i >= 0; i-- {
		if latestChanges[i].LineNum <= prLineNr {
			return latestChanges[i].LineNum, nil
		}
	}

	return 1, nil
}
