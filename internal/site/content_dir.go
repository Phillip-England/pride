package site

import "github.com/Phillip-England/pride/internal/syserr"

type ContentDir struct {
	Path string
}

func NewContentDir(path string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	return dir, nil
}