package syserr

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Default struct {
	message string
	code    ErrCode
}

func DefaultNew(code ErrCode, err error) Default {
	return Default{
		message: err.Error(),
		code:    code,
	}
}

func (err Default) GetCode() ErrCode {
	return err.code
}

func (err Default) GetMessage() string {
	return err.message
}

func (err Default) Handle() {
	fmt.Println("running the default error handler, here is your message:")
	color.Red("[ERROR_MESSAGE]: %s", err.message)
	os.Exit(1)
}
