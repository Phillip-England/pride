package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Help struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *Help) Exec(c cmd.Cmd) *syserr.Err {
	fmt.Println(`🐍 pride - simple static sites

[new]
pride new <CONTENT_TYPE> <DESTINATION>
pride new site new-site (scaffolds a new pride project)
pride new content /index.md (generate new content for your site)`)
	return nil
}
