package main

import (
	"github.com/Phillip-England/pride/internal/cmd"
)

func main() {
	cmd, err := cmd.NewCmd()
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
