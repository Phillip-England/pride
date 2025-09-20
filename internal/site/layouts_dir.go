package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type LayoutsDir struct {
	Path        string
	LayoutFiles []LayoutFile
}

func CreateLayoutsDir(path string) (LayoutsDir, error) {
	var dir LayoutsDir
	dir.Path = path

	// Create directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	// Create default layout file
	defaultLayout, err := CreateLayoutFile(filepath.Join(path, "default.html"), `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="stylesheet" href="/static/default.css">
		<title>{{ .Meta.Title }}</title>
	</head>
  	<body>
		<div id='root'>
			{{ template "header" . }}
			<main>{{ .Content }}</main>
			{{ template "footer" . }}
		</div>
  	</body>
</html>`)
	if err != nil {
		return dir, err
	}

	dir.LayoutFiles = append(dir.LayoutFiles, defaultLayout)
	return dir, nil
}

func LoadLayoutsDir(path string) (LayoutsDir, error) {
	var dir LayoutsDir
	dir.Path = path
	return dir, nil
}
