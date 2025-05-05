package format

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrintToTerminal(filename string, output string) {
	header := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(header("📄 File:"), filename)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "error") {
			color.Red("❌ %s", line)
		} else if strings.Contains(line, "warn") {
			color.Yellow("⚠️  %s", line)
		} else {
			fmt.Println(line)
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}
