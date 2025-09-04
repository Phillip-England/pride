package server

import (
	"os"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Route struct {
	Path         string
	MarkdownFile site.MarkdownFile
	LayoutPath   string
	Title        string
}

// 1. check if the .md file has a valid layout (one that exists)
func NewRoute(mdFile site.MarkdownFile) (Route, *syserr.Err) {
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
	return route, nil
}
