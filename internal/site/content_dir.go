package site

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Phillip-England/pride/internal/syserr"
)

type ContentDir struct {
	Path string
	MarkdownFiles []MarkdownFile
}

func NewContentDir(path string, configFile ConfigFile) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	indexMdFile, serr := NewMarkdownFile(filepath.Join(path, "index.md"), fmt.Sprintf(`+++
title = "Home Page"
dob = "%s"
draft = %t
template = "/templates/default.html"
+++

# Welcome
This is the home page!
`, time.Now().UTC().Format(time.RFC3339), true), configFile)
	if serr != nil {
		return dir, serr
	}
	dir.MarkdownFiles = append(dir.MarkdownFiles, indexMdFile)
	return dir, nil
}

func LoadContentDir(path string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}
	return dir, nil
}