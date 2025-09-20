package op

import (
	"os"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpNewContent struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpNewContent) Exec(c cmd.Cmd) error {
	cmdNew, ok := c.(*cmd.CmdNew)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	_, serr := OperationNewContent(cmdNew.ArgDestination, true)
	if serr != nil {
		return serr
	}
	return nil
}

func OperationNewContent(destination string, isDraft bool) (site.MarkdownFile, error) {
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return site.MarkdownFile{}, serr
	}
	// new content must be unique
	_, err := os.Stat(destination)
	if err == nil {
		return site.MarkdownFile{}, syserr.New(syserr.Here(), "%s already exists", destination)
	}
	// the names within menuNames must all begin with an uppercase letter
	mdFile, serr := site.CreateMarkdownFile(destination, "", isDraft, "/layouts/default.html", dir.ConfigFile, dir.Path, dir.ContentDir.Path)
	if serr != nil {
		return mdFile, serr
	}
	return mdFile, nil
}
