package cmd

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Help struct {
	Flag Flag
}

func HelpNew(flag Flag) (*Help, syserr.SysErr) {
	cmd := &Help{}
	cmd.Flag = flag
	return cmd, nil
}

func (cmd Help) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Help) Exec() syserr.SysErr {
	PrintHelpScreen()
	return nil
}

func PrintHelpScreen() {
	fmt.Println(`🐍 pride - simple static sites

[new]
pride new <CONTENT_TYPE> <DESTINATION>
pride new site new-site (scaffolds a new pride project)
pride new content /index.md (generate new content for your site)`)
}
