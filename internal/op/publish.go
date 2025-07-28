package op

import (
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Publish struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *Publish) Exec(c cmd.Cmd) syserr.SysErr {
	println("publishing")
	return nil
}
