package op

import (
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpNewSite struct {
	Code int
}

func (op *OpNewSite) Exec(c cmd.Cmd) *syserr.Err {
	cmdNew, ok := c.(*cmd.CmdNew)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	_, serr := site.CreatePrideDir(cmdNew.ArgDestination)
	if serr != nil {
		return serr
	}
	return nil
}
