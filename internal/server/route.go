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
	HtmlBytes    []byte
}

func NewRoute(prideDirPath string, mdFile site.MarkdownFile) (*Route, error) {
	route := &Route{}
	route.MarkdownFile = mdFile
	route.Path = mdFile.ServerPath
	route.Title = mdFile.Title
	// check if the .md file has a valid layout (one that exists)
	_, err := os.Stat(mdFile.LayoutPath)
	if err != nil {
		return route, syserr.New(syserr.Here(), ".md file at %s points to a layout at %s which does not exist", mdFile.Path, mdFile.LayoutPath)
	}
	route.LayoutPath = mdFile.LayoutPath
	// get the layout name
	route.LayoutName = filepath.Base(route.LayoutPath)
	// get the relative path
	relativePath, err := filepath.Rel(prideDirPath, mdFile.Path)
	if err != nil {
		return route, syserr.New(syserr.Here(), "%s", err.Error())
	}
	route.RelativePath = relativePath
	return route, nil
}
