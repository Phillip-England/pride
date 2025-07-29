package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Template struct {
	Path string
	Text string
}

func TemplateNew(path string) Template {
	var config Template
	config.Path = path
	config.Text = `{{ define "header" }}
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
{{ end }}

`
	return config
}

func (f Template) Create() *syserr.Err {
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.Here(), "unanticipated error when creating %s", f.Path)
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return syserr.New(syserr.Here(), "unanticipated error when writing to %s", f.Path)
	}
	return nil
}
