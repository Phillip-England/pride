package server

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

func LoadLayouts(t *template.Template, dir string) (*template.Template, error) {
	var potErr error
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
			return errors.New("")
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			_, err := t.ParseFiles(path)
			if err != nil {
				potErr = syserr.New(syserr.Here(), "%s", err.Error())
				return errors.New("")
			}
		}
		return nil
	})
	return t, potErr
}

func LoadTemplates(t *template.Template, dir string) (*template.Template, error) {
	var potErr error
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
			return errors.New("")
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			_, err := t.ParseFiles(path)
			if err != nil {
				potErr = syserr.New(syserr.Here(), "%s", err.Error())
				return errors.New("")
			}
		}
		return nil
	})
	return t, potErr
}
