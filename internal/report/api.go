package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

const (
	// DefaultAPIEndpointEnvVar is the environment variable name for the API endpoint
	DefaultAPIEndpointEnvVar = "QREVIEW_API_ENDPOINT"
)

// NewAPI creates an API reporter that sends reports to an API endpoint
func newAPI(path, reportName string) Reporter {
	return &apiReporter{
		path:           path,
		reportName:     reportName,
		processedFiles: []string{},
	}
}

type apiReporter struct {
	path           string
	reportName     string
	processedFiles []string
}

type apiPayload struct {
	FileName string `json:"fileName"`
	Content  string `json:"content"`
	Title    string `json:"title"`
}

// Report implements Reporter.
func (a *apiReporter) Report(fileName, mdContent string) error {
	// Convert markdown to HTML
	markdown := []byte(mdContent)
	var buf bytes.Buffer
	if err := goldmark.Convert(markdown, &buf); err != nil {
		return fmt.Errorf("could not convert file %w", err)
	}

	// Get the API endpoint from environment variable
	apiEndpoint := os.Getenv(DefaultAPIEndpointEnvVar)
	if apiEndpoint == "" {
		return fmt.Errorf("API endpoint environment variable %s is not set", DefaultAPIEndpointEnvVar)
	}

	// Prepare the payload
	payload := apiPayload{
		FileName: a.getRelPath(fileName),
		Content:  buf.String(),
		Title:    "Code review Report",
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Send the POST request
	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Add to processed files
	a.processedFiles = append(a.processedFiles, a.getRelPath(fileName))

	return nil
}

// Summary implements Reporter.
func (a *apiReporter) Summary(fileName string) error {
	// Get the API endpoint from environment variable
	apiEndpoint := os.Getenv(DefaultAPIEndpointEnvVar)
	if apiEndpoint == "" {
		return fmt.Errorf("API endpoint environment variable %s is not set", DefaultAPIEndpointEnvVar)
	}

	// Create summary content with links to all processed files
	var summaryContent strings.Builder
	summaryContent.WriteString("<h1>Code Review Report</h1>\n<ul>\n")
	for _, file := range a.processedFiles {
		title := strings.TrimSuffix(file, ".html")
		summaryContent.WriteString(fmt.Sprintf("  <li><a href=\"%s\">%s</a></li>\n", file, title))
	}
	summaryContent.WriteString("</ul>")

	// Prepare the payload
	payload := apiPayload{
		FileName: a.getRelPath(fileName),
		Content:  summaryContent.String(),
		Title:    "Code Review Summary",
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Send the POST request
	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *apiReporter) getRootReportPath() string {
	rootPath := a.path
	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}

	reportName := a.reportName
	if !strings.HasSuffix(reportName, "/") {
		reportName += "/"
	}

	return rootPath + reportName
}

func (a *apiReporter) getRelPath(fileName string) string {
	return fileName + ".html"
}