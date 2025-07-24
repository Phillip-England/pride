package syserr

type ErrCode int

const (
	CodeHelp ErrCode = iota
	CodeDefault
	CodeDev
	CodeMia
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
	default:
		return DefaultNew(code, err)
	}
}
