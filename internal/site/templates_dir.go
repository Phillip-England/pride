package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type TemplatesDir struct {
	Path string
	TemplateFiles []TemplateFile
}

func NewTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defaultTemplate, serr := NewTemplateFile(filepath.Join(path, "default.html"), `{{ define "header" }}
	<header>
		<h1>{{ .Pride.Ctx.Title }}</h1>
	</header>
{{ end }}
	
{{ define "nav" }}
	{{ .Pride.Nav.main }}
{{ end }}

{{ define "footer" }}
	<footer>
		<p>Created: {{ .Pride.Ctx.dob }}</p>
	</footer>
{{ end }}`)
	if serr != nil {
		return dir, serr
	}
	dir.TemplateFiles = append(dir.TemplateFiles, defaultTemplate)
	return dir, nil
}

func LoadTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	return dir, nil
}