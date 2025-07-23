package pride

import (
	"fmt"
	"os"
	"strings"
)

type CmdMake struct {
	Flag           CmdFlag
	ArgContentType string
	ArgDestination string
	ContentType    contentType
}

func CmdMakeNew(flag CmdFlag) (*CmdMake, SysErr) {
	cmd := &CmdMake{}
	argContentType, err := GetArg(2)
	if err != nil {
		return cmd, SysErrNew(SysErrCodeHelp, fmt.Errorf("missing <CONTENT_TYPE> and <DESTINATION> in 'pride make'"))
	}
	contentType, err := contentTypeNew(argContentType)
	if err != nil {
		return cmd, SysErrNew(SysErrCodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride make'"))
	}
	cmd.Flag = flag
	cmd.ArgContentType = argContentType
	cmd.ContentType = contentType
	argDestination, err := GetArg(3)
	if err != nil {
		return cmd, SysErrNew(SysErrCodeHelp, fmt.Errorf("missing <DESTINATION> in 'pride make'"))
	}
	cmd.ArgDestination = argDestination
	return cmd, nil
}

func (cmd CmdMake) GetFlag() CmdFlag {
	return cmd.Flag
}

func (cmd CmdMake) Exec() SysErr {
	switch cmd.ContentType {
	case contentTypeSite:
		op, err := OpNew(OpCodeMakeSite, cmd)
		if err != nil {
			return err
		}
		err = op.Run(cmd)
		if err != nil {
			return err
		}
		return nil
	default:
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("invalid <CONTENT_TYPE> passed to 'pride make'"))
	}
}

type contentType int

const (
	contentTypeSite contentType = iota
)

func contentTypeNew(contentType string) (contentType, error) {
	switch contentType {
	case "site":
		return contentTypeSite, nil
	default:
		return 0, fmt.Errorf("invalid content type provided")
	}
}

func makeSite(dest string) SysErr {
	fmt.Printf("making new site in %s..\n", dest)
	err := os.Mkdir(dest, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> %s already exists, please provide a new name or delete %s", dest, dest))
		}
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> unspecified error when creating %s, are you using a valid directory name for your site?", dest))

	}
	return nil
}
