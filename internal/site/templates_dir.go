package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type TemplatesDir struct {
	Path string
}

func CreateTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	return dir, nil
}

func LoadTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	return dir, nil
}
