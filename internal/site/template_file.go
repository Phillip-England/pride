package site

import (
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type TemplateFile struct {
	Path            string
	Text string
}

func NewTemplateFile(path string, text string) (TemplateFile, *syserr.Err) {
	var templateFile TemplateFile
	file, err := os.Create(path)
	if err != nil {
		return templateFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	file.WriteString(text)
	templateFile.Text = text
	return templateFile, nil
}
