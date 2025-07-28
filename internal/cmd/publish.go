package cmd

import (
	"fmt"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Publish struct {
	Flag           Flag
	ArgContentPath string
}

func PublishNew(flag Flag) (*Publish, syserr.SysErr) {
	cmd := &Publish{}
	cmd.Flag = flag
	argContentPath, err := GetArg(2)
	if err != nil {
		return cmd, syserr.HelpNew(fmt.Errorf(`missing <CONTENT-PATH> in 'pride publish'`))
	}
	argContentPath = strings.TrimPrefix(argContentPath, ".")
	argContentPath = strings.TrimPrefix(argContentPath, "./")
	// argContentPath = strings.TrimPrefix(argContentPath, "content")
	return cmd, nil
}

func (cmd Publish) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Publish) GetOpCode() (int, syserr.SysErr) {
	return 4, nil
}
