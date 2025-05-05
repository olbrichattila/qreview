package report

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func newMarkdown(reportName string) Reporter {
	return &mdReporter{reportName: reportName}
}

type mdReporter struct {
	reportName string
}

// Report implements Reporter.
func (m *mdReporter) Report(_, mdContent string) error {
	m.displayMd(mdContent)
	return nil
}

// Summary implements Reporter.
func (m *mdReporter) Summary(_ string) error {
	// We do not summarize on screen
	return nil
}

func (c *mdReporter) displayMd(content string) {
	fmt.Println(c.reportName)
	out, err := glamour.Render(string(content), "dark")
	if err != nil {
		fmt.Println(content)
	}

	fmt.Print(out)
}
