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

func (err ErrHelp) Handle() {
	PrintHelpScreen()
	fmt.Print("\n")
	fmt.Println("ops! we encountered an error while running pride, here is the error:")
	color.Red(err.message)
	os.Exit(1)
}
