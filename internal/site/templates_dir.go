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

func CreateTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	defaultTemplates, serr := CreateTemplateFile(filepath.Join(path, "templates.html"), `{{ define "header"}}
	<header>some header</header>
{{ end }}

{{ define "footer" }}
	<footer>some footer</footer>
{{ end }}
	`)
	if serr != nil {
		return dir, serr
	}
	dir.TemplateFiles = append(dir.TemplateFiles, defaultTemplates)
	return dir, nil
}

func LoadTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	return dir, nil
}
