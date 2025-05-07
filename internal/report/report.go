// Report creates a report from the answer to the provided formats
package report

import "fmt"

type Kind string

const (
	KindHTML     Kind = "html"
	KindMarkdown Kind = "markdown"
	KindSave     Kind = "save"
	KindAPI      Kind = "api"
)

type Reporter interface {
	Report(fileName, mdContent string) error
	Summary(fileName string) error
}

func New(rType Kind, path, reportName string) (Reporter, error) {
	switch rType {
	case KindHTML:
		return newHTML(path, reportName), nil
	case KindMarkdown:
		return newMarkdown(reportName), nil
	case KindSave:
		return newSaveResponse(path, reportName), nil
	case KindAPI:
		return newAPI(path, reportName), nil
	default:
		return nil, fmt.Errorf("invalid report type %s", rType)
	}
}
