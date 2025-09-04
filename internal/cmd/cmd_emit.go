package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdEmit struct {
	ArgEmitKey string
	EmitKey    emitKey
}

func NewCmdEmit() (*CmdEmit, *syserr.Err) {
	cmd := &CmdEmit{}
	argEmitKey, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <EMITKEY> in `pride emit`")
	}
	emitKey, serr := validateEmitKey(argEmitKey)
	if serr != nil {
		return cmd, serr
	}
	cmd.EmitKey = emitKey
	return cmd, nil
}

func (cmd CmdEmit) GetOpCode() (int, *syserr.Err) {
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

func validateEmitKey(argEmitKey string) (emitKey, *syserr.Err) {
	switch argEmitKey {
	case "nav":
		return emitKeyNav, nil
	default:
		return -1, syserr.New(syserr.Here(), "invalid <EMITKEY> in `pride emit`")
	}

}
