package cmd

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type Cmd interface {
	GetOpCode() (int, *syserr.Err)
}

func Extract() (Cmd, *syserr.Err) {
	id := NewCmdId()
	switch id {
	case FlagNewNew:
		cmd, err := NewCmdNew()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagEmit:
		cmd, serr := NewCmdEmit()
		if serr != nil {
			return nil, serr
		}
		return cmd, nil
	case FlagServe:
		cmd, serr := NewCmdServe()
		if serr != nil {
			return nil, serr
		}
		return cmd, nil
	default:
		cmd, err := NewCmdHelp()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
}
