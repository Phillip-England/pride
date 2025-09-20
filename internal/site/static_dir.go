package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type StaticDir struct {
	Path     string
	CssFiles []CssFile
}

func CreateStaticDir(path string) (StaticDir, error) {
	var dir StaticDir
	dir.Path = path

	if err := os.MkdirAll(path, 0755); err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	defaultCssFile, err := CreateCssFile(filepath.Join(path, "default.css"), `* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: sans-serif;
  line-height: 1.5;
  background: #fff;
  color: #000;
}`)
	if err != nil {
		return dir, err
	}

	dir.CssFiles = append(dir.CssFiles, defaultCssFile)
	return dir, nil
}

func LoadStaticDir(path string) (StaticDir, error) {
	var dir StaticDir
	dir.Path = path
	return dir, nil
}
