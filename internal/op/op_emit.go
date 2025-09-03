package op

import (
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpEmit struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpEmit) Exec(c cmd.Cmd) *syserr.Err {
	// now we can emit..!
	return nil
}
