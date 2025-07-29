package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type NewSite struct {
	Code int
}

func (op *NewSite) Exec(c cmd.Cmd) *syserr.Err {
	cmdNew, ok := c.(*cmd.New)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	fmt.Printf("🧬 making new site at %s\n", cmdNew.ArgDestination)
	dir := site.RootNew(cmdNew.ArgDestination)
	err := dir.Create()
	if err != nil {
		return err
	}
	return nil
}
