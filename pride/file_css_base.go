package pride

import (
	"fmt"
	"os"
)

type FileCssBase struct {
	Path string
	Text string
}

func FileCssBaseNew(path string) FileCssBase {
	var config FileCssBase
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

func (f FileCssBase) Create() SysErr {
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
