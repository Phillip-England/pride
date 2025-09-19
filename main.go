package main

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/op"
)

func main() {
	command, err := cmd.Extract()
	if err != nil {
		handleErr(err)
		return
	}
	code, err := command.GetOpCode()
	if err != nil {
		handleErr(err)
		return
	}
	operation, err := op.OpNew(code)
	if err != nil {
		handleErr(err)
		return
	}
	err = operation.Exec(command)
	if err != nil {
		handleErr(err)
		return
	}
}

func handleErr(err error) {
	fmt.Println(err.Error())
}
