package cmd

import (
	"strconv"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdServe struct {
	ArgPort string
	Port    int
}

func NewCmdServe() (*CmdServe, *syserr.Err) {
	cmd := &CmdServe{}
	argPort, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <PORT> in `pride serve`")
	}
	port, err := strconv.Atoi(argPort)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "%s was provided for <PORT> in `pride serve` but only numbers are allowed", argPort)
	}
	if port > 65000 || port < 0 {
		return cmd, syserr.New(syserr.Here(), "<PORT> used in `pride serve` must be between 0 and 65000")
	}
	cmd.Port = port
	return cmd, nil
}

func (cmd CmdServe) GetOpCode() (int, *syserr.Err) {
	return 4, nil
}
