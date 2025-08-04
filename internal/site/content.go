package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Content struct {
	Path    string
	Text    string
	Dob     string
	IsDraft bool
	Title   string
}

func ContentNew(path string, title string, isDraft bool) *Content {
	var f Content
	f.Title = title
	f.Path = path
	f.Dob = time.Now().UTC().Format(time.RFC3339)
	f.IsDraft = isDraft
	f.Text = fmt.Sprintf(`+++
title = "%s"
dob = "%s"
draft = %t
template = "/templates/default.html"
+++

# Welcome
This is the home page!
`, f.Title, f.Dob, f.IsDraft)
	return &f
}

func (f Content) Create() *syserr.Err {
	dir := filepath.Dir(f.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return syserr.New(syserr.Here(), "failed to create parent directories for %s => %s", f.Path, err.Error())
	}
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.Here(), "unanticipated error when creating %s => %s", f.Path, err.Error())
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return syserr.New(syserr.Here(), "unanticipated error when writing to %s => %s", f.Path, err.Error())
	}
	return nil
}

func GetContentPaths(contentDir string) ([]string, *syserr.Err) {
	paths := []string{}
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return paths, syserr.New(syserr.Here(), "%s", err.Error())
	}
	sort.Strings(paths)
	return paths, nil
}

func ContentLoadAll() ([]*MarkdownFile, *syserr.Err) {
	content := []*MarkdownFile{}
	config, serr := LoadConfigFile()
	if serr != nil {
		return content, serr
	}
	paths, serr := GetContentPaths("./content")
	if serr != nil {
		return content, serr
	}
	for _, path := range paths {
		mdFile, serr := MarkdownFileLoad(path, "content", config.Theme, "./templates")
		if serr != nil {
			return content, serr
		}
		content = append(content, mdFile)
	}
	return content, nil
}