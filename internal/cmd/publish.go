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
	argContentPath, err := ArgEnforePath(2)
	if err != nil {
		return cmd, syserr.New(syserr.Here(), `<CONTENT-PATH> error\n%s`, err.Error())
	}
	config, serr := site.ConfigLoad()
	if serr != nil {
		return cmd, serr
	}
	cmd.Config = config
	cmd.ArgContentPath = argContentPath
	println(cmd.Config.Path)
	return cmd, nil
}

func (cmd Publish) GetFlag() Flag {
	return cmd.Flag
}

func (cmd Publish) GetOpCode() (int, *syserr.Err) {
	return 4, nil
}
