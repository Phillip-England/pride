package syserr

type ErrCode int

const (
	CodeHelp ErrCode = iota
	CodeDefault
	CodeDev
	CodeMia
	CodeLib
)

type SysErr interface {
	Handle()
	GetCode() ErrCode
	GetMessage() string
}

func New(code ErrCode, err error) SysErr {
	switch code {
	case CodeHelp:
		return HelpNew(code, err)
	case CodeDev:
		return DevNew(code, err)
	case CodeMia:
		return MiaNew(code, err)
	case CodeLib:
		return LibNew(code, err)
	default:
		return DefaultNew(code, err)
	}
}
