package syserr

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Help struct {
	message string
	code    ErrCode
}

func HelpNew(err error) Help {
	return Help{
		message: err.Error(),
		code:    CodeHelp,
	}
}

func (err Help) GetCode() ErrCode {
	return err.code
}

func (err Help) GetMessage() string {
	return err.message
}

func (err Help) Handle() {
	fmt.Println("ops! error encountered, run 'pride help' for more info:")
	color.Red(err.message)
	os.Exit(1)
}
