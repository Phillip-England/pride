package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CssFile struct {
	Path string
	Text string
}

func NewCssFile(path string, text string) (CssFile, *syserr.Err) {
	var cssFile CssFile
	file, err := os.Create(path)
	if err != nil {
		return cssFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	file.WriteString(text)
	cssFile.Text = text
	return cssFile, nil
}