package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpHelp struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpHelp) Exec(c cmd.Cmd) *syserr.Err {
	fmt.Println(`🐍 pride
A dead-simple, hugo-inspired, static-site generator.

[new]
- pride new site new-site
- pride new content ./content/index.md
[emit]
- pride emit nav ./content ./nav.html
- pride emit nav ./content/posts ./nav-posts.html
[serve]
- pride serve 8080`)
	return nil
}
