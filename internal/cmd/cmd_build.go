package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdBuild struct {
}

func NewCmdBuild() (*CmdBuild, *syserr.Err) {
	cmd := &CmdBuild{}
	return cmd, nil
}

func (cmd CmdBuild) GetOpCode() (int, *syserr.Err) {
	return 5, nil
}
