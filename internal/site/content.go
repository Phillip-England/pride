package site

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Content struct {
	Path string
	Text string
	Dob  string
}

func ContentNew(path string) Content {
	var f Content
	f.Path = path
	f.Dob = time.Now().UTC().Format(time.RFC3339)
	f.Text = fmt.Sprintf(`---
title: "Home Page"
dob: "%s"
draft: "true"
template: "/default_template.html"
nav: "main"
---

# Welcome
This is the home page!
`, f.Dob)
	return f
}

func (f Content) Create() syserr.SysErr {
	dir := filepath.Dir(f.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("failed to create parent directories for %s => %s", f.Path, err.Error()))
	}
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when creating %s => %s", f.Path, err.Error()))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when writing to %s => %s", f.Path, err.Error()))
	}
	return nil
}
