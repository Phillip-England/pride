package cmd

import (
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdBuild struct {
	ArgDest string
}

func NewCmdBuild() (*CmdBuild, *syserr.Err) {
	cmd := &CmdBuild{}
	argDest, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "<DESTINATION> missing from 'pride build'")
	}
	argDest, err = filepath.Abs(argDest)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "<DESTINATION> could not be converted into an absolute filepath in 'pride build'")
	}
	cmd.ArgDest = argDest
	return cmd, nil
}

func (cmd CmdBuild) GetOpCode() (int, *syserr.Err) {
	return 5, nil
}
