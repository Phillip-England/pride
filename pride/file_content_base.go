package pride

import (
	"fmt"
	"os"
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
template: "/default_template.html"
nav: "main"
---

# Welcome
This is the home page!
`, f.Dob)
	return f
}

func (f FileContentBase) Create() SysErr {
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when creating %s", f.Path))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when writing to %s", f.Path))
	}
	return nil
}
