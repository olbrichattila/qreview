// Package parentsummary generates a summary for child folders
package parentsummary

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template/template.html
var templateFS embed.FS

type Item struct {
	Href  string
	Title string
}

type pageData struct {
	Title string
	Items []Item
}

func Generate(rootPath, reportPath string) error {
	// Check if reportPath exists, skip if it doesn't
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		return nil
	}

	files, err := getHtmlFiles(reportPath)
	if err != nil {
		return err
	}

	err = save(reportPath+"/index.html", files)
	if err != nil {
		return err
	}

	parent := filepath.Dir(reportPath)
	if parent != "." {
		err := Generate(rootPath, parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func getHtmlFiles(reportPath string) ([]string, error) {
	var result []string

	var walk func(string) error
	walk = func(path string) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		foundIndex := false
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() == "index.html" && path != reportPath {
				relPath := strings.TrimPrefix(path, reportPath+"/")
				result = append(result, relPath+"/index.html")
				foundIndex = true
				break
			}
		}

		if foundIndex {
			// Don't go deeper
			return nil
		}

		for _, entry := range entries {
			if entry.IsDir() {
				err := walk(filepath.Join(path, entry.Name()))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := walk(reportPath); err != nil {
		return nil, err
	}

	return result, nil
}

func save(fileName string, links []string) error {
	tmpl, err := template.ParseFS(templateFS, "template/template.html")
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {

		return err
	}
	defer file.Close()

	items := make([]Item, len(links))
	for i, link := range links {
		items[i] = Item{
			Href:  link,
			Title: strings.TrimSuffix(link, "/index.html"),
		}
	}
	pageData := pageData{
		Title: "Code review Report",
		Items: items,
	}

	err = tmpl.Execute(file, pageData)
	if err != nil {
		return err
	}

	return nil
}
