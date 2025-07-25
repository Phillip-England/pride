package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type BuildNav struct {
	Code int
}

func (op *BuildNav) Exec(c cmd.Cmd) syserr.SysErr {
	_, ok := c.(*cmd.Build)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure, did you use pointers correctly? did you return a valid op code?"))
	}
	fmt.Printf("🧬 building site navigation\n")

	return nil
}
