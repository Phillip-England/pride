package op

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Op interface {
	Exec(c cmd.Cmd) syserr.SysErr
}

func OpNew(code int) (Op, syserr.SysErr) {
	switch code {
	case 0:
		return &Help{Code: code}, nil
	case 1:
		return &NewSite{Code: code}, nil
	case 2:
		return &NewContent{Code: code}, nil
	default:
		return nil, syserr.New(syserr.CodeDev, fmt.Errorf("provided an invalid code to OpNew in ./internal/op/op.go"))
	}
}
