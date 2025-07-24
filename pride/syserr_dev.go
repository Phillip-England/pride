package pride

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ErrDev struct {
	message string
	code    ErrCode
}

func ErrDevNew(code ErrCode, err error) ErrDev {
	return ErrDev{
		message: err.Error(),
		code:    code,
	}
}

func (err ErrDev) GetCode() ErrCode {
	return err.code
}

func (err ErrDev) GetMessage() string {
	return err.message
}

func (err ErrDev) Handle() {
	fmt.Println("this error occurred because of a developer mistake:")
	color.Red(err.message)
	os.Exit(1)
}
