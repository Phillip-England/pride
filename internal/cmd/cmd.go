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
	default:
		cmd, err := NewCmdHelp()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
}
