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

func CreateContentDir(path string, configFile ConfigFile, prideDirPath string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	dir.MarkdownFiles = []MarkdownFile{}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	indexMdFile, serr := CreateMarkdownFile(filepath.Join(path, "index.md"), "Home Page", true, []string{}, "/templates/default.html", configFile, prideDirPath, dir.Path)
	if serr != nil {
		return dir, serr
	}
	dir.MarkdownFiles = append(dir.MarkdownFiles, indexMdFile)
	return dir, nil
}

func LoadContentDir(path string, theme string, prideDirPath string) (ContentDir, *syserr.Err) {
	var dir ContentDir
	dir.Path = path
	var potErr *syserr.Err
	mdFiles := []MarkdownFile{}
	err := filepath.Walk(dir.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}
		mdFile, serr := LoadMarkdownFile(path, theme, prideDirPath, dir.Path)
		if serr != nil {
			potErr = serr
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
