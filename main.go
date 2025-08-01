package main

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/site"
)

func main() {

	dir, serr := site.NewPrideDir(".")
	if serr != nil {
		serr.Print()
		return
	}
	fmt.Println(dir)

	// command, err := cmd.Extract()
	// if err != nil {
	// 	err.Print()
	// 	return
	// }
	// code, err := command.GetOpCode()
	// if err != nil {
	// 	err.Print()
	// 	return
	// }
	// operation, err := op.OpNew(code)
	// if err != nil {
	// 	err.Print()
	// 	return
	// }
	// err = operation.Exec(command)
	// if err != nil {
	// 	err.Print()
	// 	return
	// }
}
