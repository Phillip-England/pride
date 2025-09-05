package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpServe struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpServe) Exec(c cmd.Cmd) *syserr.Err {
	cmdServe, ok := c.(*cmd.CmdServe)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	port := cmdServe.Port
	serr := OperationStartServer(port)
	if serr != nil {
		return serr
	}
	return nil
}

// 1. load the project dir into memory
// 2. create a new server
func OperationStartServer(port int) *syserr.Err {
	// 1
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	// 2
	svr, serr := server.NewServer(port, dir)
	if serr != nil {
		return serr
	}
	fmt.Println(svr.Layouts)
	return nil
}
