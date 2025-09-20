package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type TemplateFile struct {
	Path string
	Text string
}

func CreateTemplateFile(path string, text string) (TemplateFile, error) {
	var f TemplateFile
	file, err := os.Create(path)
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	file.WriteString(text)
	f.Text = text
	return f, nil
}
