package main

import (
	"github.com/Phillip-England/pride/pride"
)

func main() {
	cmd, err := pride.CmdNew()
	if err != nil {
		err.Handle()
		return
	}
	err = cmd.Exec()
	if err != nil {
		err.Handle()
		return
	}
}
