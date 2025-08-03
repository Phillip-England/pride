package site

import "github.com/Phillip-England/pride/internal/syserr"

type TemplatesDir struct {
	Path string
}

func NewTemplatesDir(path string) (TemplatesDir, *syserr.Err) {
	var dir TemplatesDir
	dir.Path = path
	return dir, nil
}