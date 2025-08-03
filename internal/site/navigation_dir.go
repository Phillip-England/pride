package site

import "github.com/Phillip-England/pride/internal/syserr"

type NavigationDir struct {
	Path string
}

func NewNavigationDir(path string) (NavigationDir, *syserr.Err) {
	var dir NavigationDir
	dir.Path = path
	return dir, nil
}