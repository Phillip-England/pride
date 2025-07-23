package pride

import (
	"fmt"
	"os"
	"strings"
)

type OpMakeSite struct {
	Code OpCode
	Cmd  Cmd
}

func (op *OpMakeSite) Run(cmd Cmd) SysErr {
	cmdMake, ok := cmd.(CmdMake)
	if !ok {
		return SysErrNew(SysErrCodeDev, fmt.Errorf("type assertion failure"))
	}
	dest := cmdMake.ArgDestination
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
