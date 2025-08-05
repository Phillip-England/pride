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

func (op *OpNewContent) Exec(c cmd.Cmd) *syserr.Err {
	cmdNew, ok := c.(*cmd.CmdNew)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	isDraft := !cmd.HasFlag("-f")
	_, serr := OperationNewContent(cmdNew.ArgDestination, isDraft)
	if serr != nil {
		return serr
	}
	return nil
}

func OperationNewContent(destination string, isDraft bool) (site.MarkdownFile, *syserr.Err) {
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return site.MarkdownFile{}, serr
	}
	_, err := os.Stat(destination)
	if err == nil {
		return site.MarkdownFile{}, syserr.New(syserr.Here(), "%s already exists", destination)
	}
	mdFile, serr := site.CreateMarkdownFile(destination, "", isDraft, "/templates/default.html", dir.ConfigFile, dir.Path, dir.ContentDir.Path)
	if serr != nil {
		return mdFile, serr
	}
	return mdFile, nil
}
