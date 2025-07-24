package pride

type ErrCode int

const (
	SysErrCodeHelp ErrCode = iota
	SysErrCodeDefault
	SysErrCodeDev
	SysErrCodeMia
)

type SysErr interface {
	Handle()
	GetCode() ErrCode
	GetMessage() string
}

func SysErrNew(code ErrCode, err error) SysErr {
	switch code {
	case SysErrCodeHelp:
		return ErrHelpNew(code, err)
	case SysErrCodeDev:
		return ErrDevNew(code, err)
	case SysErrCodeMia:
		return ErrMiaNew(code, err)
	default:
		return ErrDefaultNew(code, err)
	}
}
