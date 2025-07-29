package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Layout struct {
	Path string
	Text string
}

func LayoutNew(path string) Layout {
	var config Layout
	config.Path = path
	config.Text = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="stylesheet" href="/static/default.css">
		<title>{{ .Pride.Ctx.Title }}</title>
	</head>
  	<body>
		<div id='root'>
			{{ template "header" . }}
			{{ template "nav" . }}
			<main>{{ .Pride.Content }}</main>
			{{ template "footer" . }}
		</div>
  	</body>
</html>`
	return config
}

func (f Layout) Create() *syserr.Err {
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
