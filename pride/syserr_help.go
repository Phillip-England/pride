package pride

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ErrHelp struct {
	message string
	code    ErrCode
}

func ErrHelpNew(code ErrCode, err error) ErrHelp {
	return ErrHelp{
		message: err.Error(),
		code:    code,
	}
}

func (err ErrHelp) GetCode() ErrCode {
	return err.code
}

func (err ErrHelp) GetMessage() string {
	return err.message
}

func (err ErrHelp) Handle() {
	fmt.Println("ops! error encountered, run 'pride help' for more info:")
	color.Red(err.message)
	os.Exit(1)
}
