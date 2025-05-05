package report

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

//go:embed template/summary-template.html
var summaryTemplateFS embed.FS

//go:embed template/template.html
var templateFS embed.FS

// NewHTML creates a HTML reporter
func newHTML(path, reportName string) Reporter {
	return &htmlReporter{
		path:           path,
		reportName:     reportName,
		processedFiles: []string{},
	}
}

type Item struct {
	Href  string
	Title string
}

type summaryPageData struct {
	Title string
	Items []Item
}

type pageData struct {
	Title   string
	Content string
}

type htmlReporter struct {
	path           string
	reportName     string
	processedFiles []string
}

// Report implements Reporter.
func (h *htmlReporter) Report(fileName, mdContent string) error {
	reportFileName := h.getFullPath(fileName)
	markdown := []byte(mdContent)

	var buf bytes.Buffer
	if err := goldmark.Convert(markdown, &buf); err != nil {
		return fmt.Errorf("could not convert file %w", err)
	}

	err := os.MkdirAll(filepath.Dir(reportFileName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = h.save(reportFileName, buf.String())

	if err != nil {
		return err
	}

	h.processedFiles = append(h.processedFiles, h.getRelPath(fileName))

	return nil
}

func (h *htmlReporter) getRootReportPath() string {
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

func (h *htmlReporter) getFullPath(fileName string) string {
	reportPath := h.getRootReportPath()

	return reportPath + h.getRelPath(fileName)
}

func (h *htmlReporter) getRelPath(fileName string) string {
	return fileName + ".html"
}

// Summary implements Reporter.
func (h *htmlReporter) Summary(fileName string) error {
	tmpl, err := template.ParseFS(summaryTemplateFS, "template/summary-template.html")
	if err != nil {
		return err
	}

	indexHTMLFileName := h.getFullPath(fileName)

	file, err := os.Create(indexHTMLFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	items := make([]Item, len(h.processedFiles))
	for i, processedFile := range h.processedFiles {
		items[i] = Item{
			Href:  processedFile,
			Title: strings.TrimSuffix(processedFile, ".html"),
		}

	}
	pageData := summaryPageData{
		Title: "Code review Report",
		Items: items,
	}

	err = tmpl.Execute(file, pageData)
	if err != nil {
		return err
	}

	return nil
}

func (h *htmlReporter) save(fileName string, content string) error {
	tmpl, err := template.ParseFS(templateFS, "template/template.html")
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {

		return err
	}
	defer file.Close()

	pageData := pageData{
		Title:   "Code review Report",
		Content: content,
	}

	err = tmpl.Execute(file, pageData)
	if err != nil {
		return err
	}

	return nil
}
