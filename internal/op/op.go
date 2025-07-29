package op

import (
	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Op interface {
	Exec(c cmd.Cmd) *syserr.Err
}

func OpNew(code int) (Op, *syserr.Err) {
	switch code {
	case 0:
		return &Help{Code: code}, nil
	case 1:
		return &NewSite{Code: code}, nil
	case 2:
		return &NewContent{Code: code}, nil
	case 3:
		return &BuildNav{Code: code}, nil
	case 4:
		return &Publish{Code: code}, nil
	default:
		return nil, syserr.New(syserr.Here(), "provided an invalid code to OpNew in ./internal/op/op.go")
	}
}
