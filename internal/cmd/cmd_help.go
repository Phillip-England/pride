package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdHelp struct {
}

func NewCmdHelp() (*CmdHelp, *syserr.Err) {
	cmd := &CmdHelp{}
	return cmd, nil
}

func (cmd CmdHelp) GetOpCode() (int, *syserr.Err) {
	return 0, nil
}
