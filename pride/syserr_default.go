package pride

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ErrDefault struct {
	message string
	code    ErrCode
}

func ErrDefaultNew(code ErrCode, err error) ErrDefault {
	return ErrDefault{
		message: err.Error(),
		code:    code,
	}
}

func (err ErrDefault) GetCode() ErrCode {
	return err.code
}

func (err ErrDefault) GetMessage() string {
	return err.message
}

func (err ErrDefault) Handle() {
	fmt.Println("running the default error handler, here is your message:")
	color.Red("[ERROR_MESSAGE]: %s", err.message)
	os.Exit(1)
}
