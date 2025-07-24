package pride

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileContentBase struct {
	Path string
	Text string
	Dob  string
}

func FileContentBaseNew(path string) FileContentBase {
	var f FileContentBase
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

func (f FileContentBase) Create() SysErr {
	dir := filepath.Dir(f.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("failed to create parent directories for %s => %s", f.Path, err.Error()))
	}
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when creating %s => %s", f.Path, err.Error()))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when writing to %s => %s", f.Path, err.Error()))
	}
	return nil
}
