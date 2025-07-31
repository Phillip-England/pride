package cmd

import "os"

type CmdId int

const (
	FlagHelp CmdId = iota
	FlagNewNew
	FlagBuild
	FlagPublish
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
	case "build":
		return FlagBuild
	case "publish":
		return FlagPublish
	default:
		return FlagHelp
	}
}
