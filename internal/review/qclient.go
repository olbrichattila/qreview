package review

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func AnalyzeCode(content string, filename string) string {
	prompt := "Review this Go code for bugs, performance, security issues, and suggest improvements."

	fmt.Println(content)

	cmd := exec.Command("q", "ask", "--question", prompt, "--code", content, "--source", filename)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "❌ Amazon Q error: " + err.Error()
	}

	return out.String()
}

func ContainsCritical(output string) bool {
	return strings.Contains(output, "security") || strings.Contains(output, "bug") || strings.Contains(output, "vulnerability")
}
