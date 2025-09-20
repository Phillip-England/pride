package op

import (
	"github.com/Phillip-England/pride/internal/build"
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpBuild struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpBuild) Exec(c cmd.Cmd) error {
	cmdBuild, ok := c.(*cmd.CmdBuild)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	serr := build.GenerateBuild(cmdBuild.ArgDest)
	if serr != nil {
		return serr
	}
	return nil
}
