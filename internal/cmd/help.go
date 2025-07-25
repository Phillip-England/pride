package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type Help struct {
	Flag Flag
}

func HelpNew(flag Flag) (*Help, syserr.SysErr) {
	cmd := &Help{}
	cmd.Flag = flag
	return cmd, nil
}

func (cmd Help) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Help) GetOpCode() (int, syserr.SysErr) {
	return 0, nil
}
