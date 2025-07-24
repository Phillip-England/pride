package site

import (
	"fmt"
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Css struct {
	Path string
	Text string
}

func CssNew(path string) Css {
	var config Css
	config.Path = path
	config.Text = `* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: sans-serif;
  line-height: 1.5;
  background: #fff;
  color: #000;
}
`
	return config
}

func (f Css) Create() syserr.SysErr {
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when creating %s", f.Path))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when writing to %s", f.Path))
	}
	return nil
}
