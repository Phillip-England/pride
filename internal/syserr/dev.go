package syserr

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Dev struct {
	message string
	code    ErrCode
}

func DevNew(err error) Dev {
	return Dev{
		message: err.Error(),
		code:    CodeDev,
	}
}

func (err Dev) GetCode() ErrCode {
	return err.code
}

func (err Dev) GetMessage() string {
	return err.message
}

func (err Dev) Handle() {
	fmt.Println("this error occurred because of a developer mistake:")
	color.Red(err.message)
	os.Exit(1)
}
