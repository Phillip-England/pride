package pride

import (
	"fmt"
	"os"
)

type Cmd interface {
	GetFlag() CmdFlag
	Exec() SysErr
}

func CmdNew() (Cmd, SysErr) {
	flag := CliFlagNew()
	switch flag {
	case CmdFlagMake:
		cmd, err := CmdMakeNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		cmd, err := CmdHelpNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
}

func GetArg(n int) (string, error) {
	if len(os.Args) <= n {
		return "", fmt.Errorf(`out of bounds access on os.Args`)
	}
	return os.Args[n], nil
}
