package pride

type ErrCode int

const (
	SysErrCodeHelp ErrCode = iota
	SysErrCodeDefault
	SysErrCodeDev
)

type SysErr interface {
	Handle()
	GetCode() ErrCode
}

func SysErrNew(code ErrCode, err error) SysErr {
	switch code {
	case SysErrCodeHelp:
		return ErrHelpNew(code, err)
	case SysErrCodeDev:
		return ErrDevNew(code, err)
	default:
		return ErrDefaultNew(code, err)
	}
}
