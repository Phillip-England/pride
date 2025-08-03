package site

import "github.com/Phillip-England/pride/internal/syserr"

type StaticDir struct {
	Path string
}

func NewStaticDir(path string) (StaticDir, *syserr.Err) {
	var dir StaticDir
	dir.Path = path
	return dir, nil
}