package pride

import (
	"fmt"
)

type CmdHelp struct {
	Flag CmdFlag
}

func CmdHelpNew(flag CmdFlag) (*CmdHelp, SysErr) {
	cmd := &CmdHelp{}
	cmd.Flag = flag
	return cmd, nil
}

func (cmd CmdHelp) GetFlag() CmdFlag {
	return cmd.Flag
}

func (cmd CmdHelp) Exec() SysErr {
	PrintHelpScreen()
	return nil
}

func PrintHelpScreen() {
	fmt.Println(`🐍 pride - simple static sites

[make]
pride make <CONTENT-TYPE> <DESTINATION>
pride make site ./new-site (scaffolds a new pride project in <DESTINATION>)`)
}
