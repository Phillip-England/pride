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

func OperationStartServer(port int) *syserr.Err {
	fmt.Printf("serving on port %d\n", port)
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	routes := []server.Route{}
	for _, mdFile := range dir.ContentDir.MarkdownFiles {
		route, serr := server.NewRoute(mdFile)
		if serr != nil {
			return serr
		}
		routes = append(routes, route)
	}
	for _, route := range routes {
		fmt.Println(route.Path)
	}
	return nil
}
