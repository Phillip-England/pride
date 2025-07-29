package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Publish struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *Publish) Exec(c cmd.Cmd) *syserr.Err {
	cmd, ok := c.(*cmd.Publish)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	fmt.Println(cmd.ArgContentPath)
	return nil
}
