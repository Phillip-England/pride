package pride

import (
	"fmt"
	"os"
	"strings"
)

type OpMakeContent struct {
	Code OpCode
	Cmd  Cmd
}

func (op *OpMakeContent) Run(cmd Cmd) SysErr {
	cmdMake, ok := cmd.(CmdMake)
	if !ok {
		return SysErrNew(SysErrCodeDev, fmt.Errorf("type assertion failure"))
	}
	contentPath := "./content/" + cmdMake.ArgDestinationStripped
	fmt.Printf("🧬 making new content at %s\n", contentPath)
	_, syserr := FileConfigLoadFromCwd()
	if syserr != nil {
		return syserr
	}
	contentDirPath := "./content"
	info, err := os.Stat(contentDirPath)
	if os.IsNotExist(err) || !info.IsDir() {
		return SysErrNew(SysErrCodeMia, fmt.Errorf("could not located %s directory", contentDirPath))
	}
	if err != nil {
		return SysErrNew(SysErrCodeDev, fmt.Errorf("could not load information on %s after it has been confirmed to exist", contentDirPath))
	}
	f := FileContentBaseNew(contentPath)
	syserr = f.Create()
	if syserr != nil {
		if strings.Contains(syserr.GetMessage(), "file exists") {
			return SysErrNew(SysErrCodeHelp, fmt.Errorf("%s already exists", contentPath))
		}
		return syserr
	}
	return nil
}
