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
	dir := PrideDirNew(cmdMake.ArgDestination)
	err := dir.Create()
	if err != nil {
		return err
	}
	return nil
}

type PrideDir struct {
	PathRoot    string
	PathContent string
}

func PrideDirNew(path string) PrideDir {
	var dir PrideDir
	dir.PathRoot = path
	dir.PathContent = path + "/content"
	return dir
}

func (dir PrideDir) Create() SysErr {
	fmt.Printf("making new site in %s..\n", dir.PathRoot)
	err := os.Mkdir(dir.PathRoot, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> %s already exists, please provide a new name or delete %s", dir.PathRoot, dir.PathRoot))
		}
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> unspecified error when creating %s, are you using a valid directory name for your site?", dir.PathRoot))

	}
	err = os.Mkdir(dir.PathContent, 0755)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> unspecified error when creating %s, are you using a valid directory name for your site?", dir.PathContent))

	}
	return nil
}
