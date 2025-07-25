package main

import (
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/op"
)

func main() {
	command, err := cmd.CmdNew()
	if err != nil {
		err.Handle()
		return
	}
	code, err := command.GetOpCode()
	if err != nil {
		err.Handle()
		return
	}
	operation, err := op.OpNew(code)
	if err != nil {
		err.Handle()
		return
	}
	err = operation.Exec(command)
	if err != nil {
		err.Handle()
		return
	}
}
