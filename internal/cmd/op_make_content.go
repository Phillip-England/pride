package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type MakeContent struct {
	Code Code
	Cmd  Cmd
}

func (op *MakeContent) Run(c Cmd) syserr.SysErr {
	cmdMake, ok := c.(New)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure"))
	}
	contentPath := "./content/" + cmdMake.ArgDestinationStripped
	fmt.Printf("🧬 making new content at %s\n", contentPath)
	_, serr := site.ConfigLoadFromCwd()
	if serr != nil {
		return serr
	}
	contentDirPath := "./content"
	info, err := os.Stat(contentDirPath)
	if os.IsNotExist(err) || !info.IsDir() {
		return syserr.New(syserr.CodeMia, fmt.Errorf("could not located %s directory", contentDirPath))
	}
	if err != nil {
		return syserr.New(syserr.CodeDev, fmt.Errorf("could not load information on %s after it has been confirmed to exist", contentDirPath))
	}
	f := site.ContentNew(contentPath)
	serr = f.Create()
	if serr != nil {
		if strings.Contains(serr.GetMessage(), "file exists") {
			return syserr.New(syserr.CodeHelp, fmt.Errorf("%s already exists", contentPath))
		}
		return serr
	}
	return nil
}
