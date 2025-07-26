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
	ContentType            contentType
}

func NewNew(flag Flag) (*New, syserr.SysErr) {
	cmd := &New{}
	argMakeType, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("missing <CONTENT_TYPE> and <DESTINATION> in 'pride new'"))
	}
	contentType, err := contentTypeNew(argMakeType)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride new'"))
	}
	cmd.Flag = flag
	cmd.ArgMakeType = argMakeType
	cmd.ContentType = contentType
	argDestination, err := GetArg(3)
	if err != nil {
		return cmd, syserr.New(syserr.CodeHelp, fmt.Errorf("missing <DESTINATION> in 'pride new'"))
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

func (cmd New) GetOpCode() (int, syserr.SysErr) {
	switch cmd.ContentType {
	case contentTypeSite:
		return 1, nil
	case contentTypeContent:
		return 2, nil
	default:
		return -1, syserr.New(syserr.CodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride new'"))
	}
}

type contentType int

const (
	contentTypeSite contentType = iota
	contentTypeContent
)

func contentTypeNew(contentType string) (contentType, error) {
	switch contentType {
	case "site":
		return contentTypeSite, nil
	case "content":
		return contentTypeContent, nil
	default:
		return 0, fmt.Errorf("invalid content type provided")
	}
}
