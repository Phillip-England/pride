package cmd

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Build struct {
	Flag                   Flag
	ArgMakeType            string
	ArgDestination         string
	ArgDestinationStripped string
	BuildType              buildType
}

func BuildNew(flag Flag) (*Build, syserr.SysErr) {
	cmd := &Build{}

	return cmd, nil
}

func (cmd Build) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Build) GetOpCode() (int, syserr.SysErr) {
	switch cmd.BuildType {
	case buildTypeNav:
		return 3, nil
	default:
		return -1, syserr.New(syserr.CodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride new'"))
	}
}

type buildType int

const (
	buildTypeNav buildType = iota
)

func buildTypeNew(buildType string) (buildType, error) {
	switch buildType {
	case "nav":
		return buildTypeNav, nil
	default:
		return 0, fmt.Errorf("invalid build type provided")
	}
}
