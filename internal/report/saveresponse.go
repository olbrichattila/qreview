package report

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type saveReporter struct {
	path           string
	reportName     string
	processedFiles []string
}

// NewHTML creates a HTML reporter
func newSaveResponse(path, reportName string) Reporter {
	return &saveReporter{
		path:           path,
		reportName:     reportName,
		processedFiles: []string{},
	}
}

// Report implements Reporter.
func (h *saveReporter) Report(fileName, mdContent string) error {
	reportFileName := h.getFullPath(fileName)

	err := os.MkdirAll(filepath.Dir(reportFileName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = h.save(reportFileName, mdContent)

	if err != nil {
		return err
	}

	return nil
}

// Some code duplication here with html, use composition or something instead
func (h *saveReporter) getRootReportPath() string {
	rootPath := h.path
	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}

	reportName := h.reportName
	if !strings.HasSuffix(reportName, "/") {
		reportName += "/"
	}

	return rootPath + reportName
}

func (h *saveReporter) getFullPath(fileName string) string {
	reportPath := h.getRootReportPath()

	return reportPath + h.getRelPath(fileName)
}

func (h *saveReporter) getRelPath(fileName string) string {
	return fileName + ".md"
}

// Summary implements Reporter.
func (h *saveReporter) Summary(fileName string) error {
	// this is not applicable for this type of reporter
	return nil
}

func (h *saveReporter) save(fileName string, content string) error {
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("could not save md file %s, %w", fileName, err)
	}

	return nil
}
