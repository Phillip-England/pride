package cmd

import (
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type CmdPublish struct {
	ArgContentPath string
	Config         site.ConfigFile
	MdFile         *site.MarkdownFile
}

func NewCmdPublish() (*CmdPublish, *syserr.Err) {
	cmd := &CmdPublish{}
	argContentPath, exists := ArgIsFilePath(2)
	if !exists {
		return cmd, syserr.New(syserr.Here(), "%s is not a valid file path", argContentPath)
	}
	config, serr := site.LoadConfigFile()
	if serr != nil {
		return cmd, serr
	}
	cmd.Config = config
	cmd.ArgContentPath = argContentPath
	return cmd, nil
}

func (cmd CmdPublish) GetOpCode() (int, *syserr.Err) {
	return 4, nil
}
