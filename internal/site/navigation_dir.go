package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type NavigationDir struct {
	Path string
}

func CreateNavigationDir(path string) (NavigationDir, *syserr.Err) {
	var dir NavigationDir
	dir.Path = path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	return dir, nil
}

func LoadNavigationDir(path string) (NavigationDir, *syserr.Err) {
	var dir NavigationDir
	dir.Path = path
	return dir, nil
}
