package server

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Route struct {
	Path         string
	MarkdownFile site.MarkdownFile
	LayoutPath   string
	LayoutName   string
	Title        string
	RelativePath string
}

// 1. check if the .md file has a valid layout (one that exists)
// 2. get the layout name
// 3. get the relative path
func NewRoute(prideDirPath string, mdFile site.MarkdownFile) (Route, *syserr.Err) {
	var route Route
	route.MarkdownFile = mdFile
	route.Path = mdFile.ServerPath
	route.Title = mdFile.Title
	// 1
	_, err := os.Stat(mdFile.LayoutPath)
	if err != nil {
		return route, syserr.New(syserr.Here(), ".md file at %s points to a layout at %s which does not exist", mdFile.Path, mdFile.LayoutPath)
	}
	route.LayoutPath = mdFile.LayoutPath
	// 2
	route.LayoutName = filepath.Base(route.LayoutPath)
	// 3.
	relativePath, err := filepath.Rel(prideDirPath, mdFile.Path)
	if err != nil {
		return route, syserr.New(syserr.Here(), "%s", err.Error())
	}
	route.RelativePath = relativePath
	return route, nil
}
