package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type TemplatesDir struct {
	Path          string
	TemplateFiles []TemplateFile
}

func CreateTemplatesDir(path string) (TemplatesDir, error) {
	var dir TemplatesDir
	dir.Path = path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	defaultTemplates, err := CreateTemplateFile(filepath.Join(path, "templates.html"), `{{ define "header"}}
	<header>some header</header>
{{ end }}

{{ define "footer" }}
	<footer>some footer</footer>
{{ end }}
	`)
	if err != nil {
		return dir, err
	}
	dir.TemplateFiles = append(dir.TemplateFiles, defaultTemplates)
	return dir, nil
}

func LoadTemplatesDir(path string) (TemplatesDir) {
	var dir TemplatesDir
	dir.Path = path
	return dir
}
