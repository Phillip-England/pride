package site

import "github.com/Phillip-England/pride/internal/syserr"

type LayoutsDir struct {
	Path string
}

func NewLayoutsDir(path string) (LayoutsDir, *syserr.Err) {
	var dir LayoutsDir
	dir.Path = path
	return dir, nil
}

func LoadLayoutsDir(path string) (LayoutsDir, *syserr.Err) {
	var dir LayoutsDir
	dir.Path = path
	return dir, nil
}