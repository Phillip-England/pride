package pride

import "fmt"

type OpCode int

const (
	OpCodeMakeSite OpCode = iota
	OpCodeMakeContent
)

type Op interface {
	Run(Cmd) SysErr
}

func OpNew(code OpCode, cmd Cmd) (Op, SysErr) {
	switch code {
	case OpCodeMakeSite:
		return &OpMakeSite{Code: code, Cmd: cmd}, nil
	case OpCodeMakeContent:
		return &OpMakeContent{Code: code, Cmd: cmd}, nil
	default:
		return nil, SysErrNew(SysErrCodeDev, fmt.Errorf("provided an invalid OpCode to OpNew"))
	}
}
