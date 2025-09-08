package cmd

import "os"

type CmdId int

const (
	FlagHelp CmdId = iota
	FlagNewNew
	FlagEmit
	FlagServe
	FlagBuild
)

func NewCmdId() CmdId {
	args := os.Args
	if len(args) <= 1 {
		return FlagHelp
	}
	firstArg := args[1]
	switch firstArg {
	case "new":
		return FlagNewNew
	case "emit":
		return FlagEmit
	case "serve":
		return FlagServe
	case "build":
		return FlagBuild
	default:
		return FlagHelp
	}
}
