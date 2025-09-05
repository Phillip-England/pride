package server

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

func LoadLayouts(dir string) (*template.Template, error) {
	tmpl := template.New("")
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

type TemplateData struct {
	Title   string
	Content string
}
