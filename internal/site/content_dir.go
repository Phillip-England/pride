package site

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type ContentDir struct {
	Path          string
	MarkdownFiles []MarkdownFile
}

func CreateContentDir(path string, configFile ConfigFile, prideDirPath string) (ContentDir, error) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}

	// Create the content directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	// Create the index.md file
	indexMdFile, err := CreateMarkdownFile(
		filepath.Join(path, "index.md"),
		"Home Page",
		true,
		"/layouts/default.html",
		configFile,
		prideDirPath,
		dir.Path,
	)
	if err != nil {
		return dir, err
	}

	dir.MarkdownFiles = append(dir.MarkdownFiles, indexMdFile)
	return dir, nil
}

func LoadContentDir(path string, theme string, prideDirPath string) (ContentDir, error) {
	var dir ContentDir
	dir.Path = path
	var potErr error
	mdFiles := []MarkdownFile{}

	// Walk through the directory and load all Markdown files
	err := filepath.Walk(dir.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		mdFile, err := LoadMarkdownFile(path, theme, prideDirPath, dir.Path)
		if err != nil {
			potErr = err
		}
		mdFiles = append(mdFiles, mdFile)
		return nil
	})
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	if potErr != nil {
		return dir, potErr
	}

	dir.MarkdownFiles = mdFiles
	return dir, nil
}
