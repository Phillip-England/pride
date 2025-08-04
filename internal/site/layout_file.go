package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type LayoutFile struct {
	Path string
	Text string
}

func CreateLayoutFile(path string, text string) (LayoutFile, *syserr.Err) {
	var layoutFile LayoutFile
	file, err := os.Create(path)
	if err != nil {
		return layoutFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	file.WriteString(text)
	layoutFile.Text = text
	return layoutFile, nil
}
