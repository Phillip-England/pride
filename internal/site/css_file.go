package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CssFile struct {
	Path string
	Text string
}

func CreateCssFile(path string, text string) (CssFile, error) {
	var cssFile CssFile
	file, err := os.Create(path)
	if err != nil {
		return cssFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return cssFile, syserr.New(syserr.Here(), "%s", err.Error())
	}

	cssFile.Text = text
	cssFile.Path = path
	return cssFile, nil
}
