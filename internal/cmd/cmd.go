package cmd

import (
	"fmt"
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Flag int

const (
	FlagHelp Flag = iota
	FlagNewNew
	FlagBuild
	FlagPublish
)

func (flag Flag) String() string {
	switch flag {
	case 0:
		return "FLAGHELP"
	case 1:
		return "FLAGNEW"
	case 2:
		return "FLAGBUILD"
	case 3:
		return "FLAGPUBLISH"
	default:
		return "FLAGHELP"
	}
}

func flagNew() Flag {
	args := os.Args
	if len(args) <= 1 {
		return FlagHelp
	}
	firstArg := args[1]
	switch firstArg {
	case "new":
		return FlagNewNew
	case "build":
		return FlagBuild
	case "publish":
		return FlagPublish
	default:
		return FlagHelp
	}
}

type Cmd interface {
	GetFlag() Flag
	GetOpCode() (int, *syserr.Err)
}

func Extract() (Cmd, *syserr.Err) {
	flag := flagNew()
	switch flag {
	case FlagNewNew:
		cmd, err := NewNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagBuild:
		cmd, err := BuildNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagPublish:
		cmd, err := PublishNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		cmd, err := HelpNew(flag)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
}

func GetArg(n int) (string, error) {
	if len(os.Args) <= n {
		return "", fmt.Errorf(`out of bounds access on os.Args`)
	}
	return os.Args[n], nil
}

func ArgIsFilePath(position int) (string, bool) {
	arg, err := GetArg(position)
	if err != nil {
		return "", false
	}
	f, err := os.OpenFile(arg, os.O_CREATE, 0644)
	if err != nil {
		return "", false
	}
	f.Close()
	os.Remove(arg)
	return arg, true
}
