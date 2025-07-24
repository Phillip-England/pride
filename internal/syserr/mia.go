package syserr

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Mia struct {
	message string
	code    ErrCode
}

func MiaNew(code ErrCode, err error) Mia {
	return Mia{
		message: err.Error(),
		code:    code,
	}
}

func (err Mia) GetCode() ErrCode {
	return err.code
}

func (err Mia) GetMessage() string {
	return err.message
}

func (err Mia) Handle() {
	fmt.Println("operation only allowed within a pride project\nto make a new project, run 'pride make site new-site'\nhere is your error message:")
	color.Red("[ERROR_MESSAGE]: %s", err.message)
	os.Exit(1)
}
