package cmd

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Code int

const (
	CodeMakeSite Code = iota
	CodeMakeContent
)

type Op interface {
	Run(Cmd) syserr.SysErr
}

func OpNew(code Code, cmd Cmd) (Op, syserr.SysErr) {
	switch code {
	case CodeMakeSite:
		return &MakeSite{Code: code, Cmd: cmd}, nil
	case CodeMakeContent:
		return &MakeContent{Code: code, Cmd: cmd}, nil
	default:
		return nil, syserr.New(syserr.CodeDev, fmt.Errorf("provided an invalid code to OpNew"))
	}
}
