package cmd

import (
	"fmt"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
)

type New struct {
	Flag                   Flag
	ArgMakeType            string
	ArgDestination         string
	ArgDestinationStripped string
	MakeType               makeType
}

func NewNew(flag Flag) (*New, syserr.SysErr) {
	cmd := &New{}
	argMakeType, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("missing <CONTENT_TYPE> and <DESTINATION> in 'pride make'"))
	}
	makeType, err := makeTypeNew(argMakeType)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride make'"))
	}
	cmd.Flag = flag
	cmd.ArgMakeType = argMakeType
	cmd.MakeType = makeType
	argDestination, err := GetArg(3)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("missing <DESTINATION> in 'pride make'"))
	}
	cmd.ArgDestination = argDestination
	if strings.HasPrefix(cmd.ArgDestination, "./") {
		cmd.ArgDestination = strings.Replace(cmd.ArgDestination, "./", "/", 1)
	}
	if !strings.HasPrefix(cmd.ArgDestination, "/") {
		cmd.ArgDestination = "./" + cmd.ArgDestination
	}
	cmd.ArgDestinationStripped = strings.Replace(cmd.ArgDestination, "./", "", 1)
	return cmd, nil
}

func (cmd New) GetFlag() Flag {
	return cmd.Flag
}

func (cmd New) Exec() syserr.SysErr {
	switch cmd.MakeType {
	case makeTypeSite:
		op, err := OpNew(CodeMakeSite, cmd)
		if err != nil {
			return err
		}
		err = op.Run(cmd)
		if err != nil {
			return err
		}
		return nil
	case makeTypeContent:
		op, err := OpNew(CodeMakeContent, cmd)
		if err != nil {
			return err
		}
		err = op.Run(cmd)
		if err != nil {
			return err
		}
		return nil
	default:
		return syserr.New(syserr.CodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride make'"))
	}
}

type makeType int

const (
	makeTypeSite makeType = iota
	makeTypeContent
)

func makeTypeNew(makeType string) (makeType, error) {
	switch makeType {
	case "site":
		return makeTypeSite, nil
	case "content":
		return makeTypeContent, nil
	default:
		return 0, fmt.Errorf("invalid content type provided")
	}
}
