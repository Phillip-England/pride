package pride

import (
	"fmt"
)

type OpMakeSite struct {
	Code OpCode
	Cmd  Cmd
}

func (op *OpMakeSite) Run(cmd Cmd) SysErr {
	cmdMake, ok := cmd.(CmdMake)
	if !ok {
		return SysErrNew(SysErrCodeDev, fmt.Errorf("type assertion failure"))
	}
	dir := DirRootNew(cmdMake.ArgDestination)
	err := dir.Create()
	if err != nil {
		return err
	}
	return nil
}
