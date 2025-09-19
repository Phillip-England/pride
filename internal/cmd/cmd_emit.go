package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdEmit struct {
	ArgEmitKey string
	EmitKey    emitKey
}

func NewCmdEmit() (*CmdEmit, error) {
	cmd := &CmdEmit{}
	argEmitKey, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <EMITKEY> in `pride emit`")
	}
	emitKey, err := validateEmitKey(argEmitKey)
	if err != nil {
		return cmd, err
	}
	cmd.EmitKey = emitKey
	return cmd, nil
}

func (cmd CmdEmit) GetOpCode() (int, error) {
	switch cmd.EmitKey {
	case emitKeyNav:
		return 3, nil
	default:
		return -1, syserr.New(syserr.Here(), "invalid <EMITKEY> in `pride emit`")
	}
}

type emitKey int

const (
	emitKeyNav emitKey = iota
)

func validateEmitKey(argEmitKey string) (emitKey, error) {
	switch argEmitKey {
	case "nav":
		return emitKeyNav, nil
	default:
		return -1, syserr.New(syserr.Here(), "invalid <EMITKEY> in `pride emit`")
	}

}
