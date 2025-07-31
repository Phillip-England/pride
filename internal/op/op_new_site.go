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
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	dir := site.RootNew(cmdNew.ArgDestination)
	err := dir.Create()
	if err != nil {
		return err
	}
	return nil
}
