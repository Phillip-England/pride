package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpPublish struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpPublish) Exec(c cmd.Cmd) *syserr.Err {
	cmd, ok := c.(*cmd.CmdPublish)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	fmt.Println(cmd.ArgContentPath)
	return nil
}
