package pride

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ErrMia struct {
	message string
	code    ErrCode
}

func ErrMiaNew(code ErrCode, err error) ErrMia {
	return ErrMia{
		message: err.Error(),
		code:    code,
	}
}

func (err ErrMia) GetCode() ErrCode {
	return err.code
}

func (err ErrMia) GetMessage() string {
	return err.message
}

func (err ErrMia) Handle() {
	fmt.Println("operation only allowed within a pride project\nto make a new project, run 'pride make site new-site'\nhere is your error message:")
	color.Red("[ERROR_MESSAGE]: %s", err.message)
	os.Exit(1)
}
