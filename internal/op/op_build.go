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

func (op *OpBuild) Exec(c cmd.Cmd) *syserr.Err {
	serr := build.GenerateBuild()
	if serr != nil {
		return serr
	}
	return nil
}
