package site

import "github.com/Phillip-England/pride/internal/syserr"

type ContentDir struct {
	Path string
	MarkdownFiles []MarkdownFile
}

func NewContentDir(path string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}

	return dir, nil
}

func LoadContentDir(path string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}
	return dir, nil
}