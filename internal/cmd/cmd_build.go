package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdBuild struct {
	ArgMakeType            string
	ArgDestination         string
	ArgDestinationStripped string
	BuildType              buildType
}

func NewCmdBuild() (*CmdBuild, *syserr.Err) {
	cmd := &CmdBuild{}
	return cmd, nil
}

func (cmd CmdBuild) GetOpCode() (int, *syserr.Err) {
	switch cmd.BuildType {
	case buildTypeNav:
		return 3, nil
	default:
		return -1, syserr.New(syserr.Here(), "invalid <CONTENT_TYPE> passed to 'pride new'")
	}
}

type buildType int

const (
	buildTypeNav buildType = iota
)
