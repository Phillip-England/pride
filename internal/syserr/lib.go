package syserr

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Lib struct {
	message string
	code    ErrCode
}

func LibNew(err error) Lib {
	return Lib{
		message: err.Error(),
		code:    CodeLib,
	}
}

func (err Lib) GetCode() ErrCode {
	return err.code
}

func (err Lib) GetMessage() string {
	return err.message
}

func (err Lib) Handle() {
	fmt.Println("error derived from a third-party library:")
	color.Red("[ERROR_MESSAGE]: %s", err.message)
	os.Exit(1)
}
