package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type LayoutFile struct {
	Path string
	Text string
}

func CreateLayoutFile(path string, text string) (LayoutFile, error) {
	var layoutFile LayoutFile
	file, err := os.Create(path)
	if err != nil {
		return layoutFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return layoutFile, syserr.New(syserr.Here(), "%s", err.Error())
	}

	layoutFile.Path = path
	layoutFile.Text = text
	return layoutFile, nil
}
