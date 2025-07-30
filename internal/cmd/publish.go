package cmd

import (
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Publish struct {
	Flag           Flag
	ArgContentPath string
	Config         site.Config
	MdFile         *site.MarkdownFile
}

func PublishNew(flag Flag) (*Publish, *syserr.Err) {
	cmd := &Publish{}
	cmd.Flag = flag
	argContentPath, exists := ArgIsFilePath(2)
	if !exists {
		return cmd, syserr.New(syserr.Here(), "%s is not a valid file path", argContentPath)
	}
	config, serr := site.ConfigLoad()
	if serr != nil {
		return cmd, serr
	}
	cmd.Config = config
	cmd.ArgContentPath = argContentPath
	return cmd, nil
}

func (cmd Publish) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Publish) GetOpCode() (int, *syserr.Err) {
	return 4, nil
}
