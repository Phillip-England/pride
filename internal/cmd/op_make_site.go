package cmd

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type MakeSite struct {
	Code Code
	Cmd  Cmd
}

func (op *MakeSite) Run(c Cmd) syserr.SysErr {
	cmdMake, ok := c.(New)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure"))
	}
	fmt.Printf("🧬 making new site at %s\n", cmdMake.ArgDestination)
	dir := site.RootNew(cmdMake.ArgDestination)
	err := dir.Create()
	if err != nil {
		return err
	}
	return nil
}
