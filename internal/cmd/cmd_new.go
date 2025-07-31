package cmd

import (
	"fmt"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdNew struct {
	ArgMakeType            string
	ArgDestination         string
	ArgDestinationStripped string
	ContentType            contentType
}

func NewCmdNew() (*CmdNew, *syserr.Err) {
	cmd := &CmdNew{}
	argMakeType, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <CONTENT-TYPE> and <DESTINATION> in 'pride new'")
	}
	contentType, err := contentTypeNew(argMakeType)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "invalid <CONTENT-TYPE> passed to 'pride new'")
	}
	cmd.ArgMakeType = argMakeType
	cmd.ContentType = contentType
	argDestination, err := GetArg(3)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <DESTINATION> in 'pride new'")
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

func (cmd CmdNew) GetOpCode() (int, *syserr.Err) {
	switch cmd.ContentType {
	case contentTypeSite:
		return 1, nil
	case contentTypeContent:
		return 2, nil
	default:
		return -1, syserr.New(syserr.Here(), "invalid <CONTENT-TYPE> passed to 'pride new'")
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
