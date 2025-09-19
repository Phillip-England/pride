package cmd

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdNew struct {
	ArgContentType string
	ArgDestination string
	ContentType    contentType
}

func NewCmdNew() (*CmdNew, error) {
	cmd := &CmdNew{}
	argContentType, err := GetArg(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <CONTENT-TYPE> and <DESTINATION> in 'pride new'")
	}
	contentType, err := contentTypeNew(argContentType)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "invalid <CONTENT-TYPE> passed to 'pride new'")
	}
	cmd.ArgContentType = argContentType
	cmd.ContentType = contentType
	argDestination, err := GetArg(3)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), "missing <DESTINATION> in 'pride new'")
	}
	cmd.ArgDestination = argDestination
	return cmd, nil
}

func (cmd CmdNew) GetOpCode() (int, error) {
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
