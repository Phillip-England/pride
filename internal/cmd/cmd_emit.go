package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdEmit struct {
}

func NewCmdEmit() (*CmdEmit, *syserr.Err) {
	cmd := &CmdEmit{}
	return cmd, nil
}

func (cmd CmdEmit) GetOpCode() (int, *syserr.Err) {
	return 3, nil
}
