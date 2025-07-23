package pride

import "os"

type CmdFlag int

const (
	CmdFlagHelp CmdFlag = iota
	CmdFlagMake
)

func (flag CmdFlag) String() string {
	switch flag {
	case 0:
		return "FLAGHELP"
	case 1:
		return "FLAGMAKE"
	default:
		return "FLAGHELP"
	}
}

func CliFlagNew() CmdFlag {
	args := os.Args
	if len(args) <= 1 {
		return CmdFlagHelp
	}
	firstArg := args[1]
	switch firstArg {
	case "make":
		return CmdFlagMake
	default:
		return CmdFlagHelp
	}
}
