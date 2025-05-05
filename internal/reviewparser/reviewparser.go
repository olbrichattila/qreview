package reviewparser

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

// Response is the entire review spited, lines are separated, rest as summary
type Response struct {
	Summary string
	Lines   map[int]string
}

func Parse(mdFile string) Response {
	currentLineNr := 1

	response := Response{
		Summary: "This is an automated review",
		Lines:   map[int]string{},
	}

	scanner := bufio.NewScanner(strings.NewReader(mdFile))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		if lineNr, comment, ok := hasLineNumber(line); ok {
			currentLineNr = lineNr
			response.Lines[lineNr] += comment + "\n"
			continue
		}

		response.Lines[currentLineNr] += line + "\n"
	}

	return response
}

func hasLineNumber(str string) (int, string, bool) {
	if lineNr, result, ok := hasLineNumberByRange(str); ok {
		return lineNr, result, ok
	}

	return hasLineNumberBySingle(str)
}

func hasLineNumberByRange(str string) (int, string, bool) {
	return hasLineNumberByRegex(str, `(?i)Line:?\s*\d+-\d+[:*]`)
}

func hasLineNumberBySingle(str string) (int, string, bool) {
	return hasLineNumberByRegex(str, `(?i)Line:?\s*\d+[:*]`)
}

func hasLineNumberByRegex(str, regex string) (int, string, bool) {
	rangeRegex := regexp.MustCompile(regex)
	rangeMatch := rangeRegex.FindStringIndex(str)
	numberRegex := regexp.MustCompile(`\d+`)

	var matchText string
	var endPos int

	if rangeMatch != nil {
		matchText = str[rangeMatch[0]:rangeMatch[1]]
		endPos = rangeMatch[1]
		lineNr := 1
		numMatch := numberRegex.FindString(matchText)
		if numMatch != "" {
			if firstNumber, err := strconv.Atoi(numMatch); err == nil {
				lineNr = firstNumber
			}
		}

		// PR is 1 indexed, make sure no 0 returned
		if lineNr == 0 {
			lineNr = 1
		}

		return lineNr, str[endPos:], true
	}

	return 0, "", false
}
