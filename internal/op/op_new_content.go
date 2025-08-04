package op

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpNewContent struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpNewContent) Exec(c cmd.Cmd) *syserr.Err {
	cmdMake, ok := c.(*cmd.CmdNew)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	prideDir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	contentPath := "./content/" + cmdMake.ArgDestinationStripped
	contentDirPath := "./content"
	_, err := os.Stat(contentDirPath)
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	fileName := filepath.Base(cmdMake.ArgDestination)
	fileName = strings.TrimSuffix(fileName, ".md")
	parts := strings.Split(fileName, "_")
	for i, part := range parts {
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	title := strings.Join(parts, " ")
	_, serr = site.CreateMarkdownFile(contentPath+".md", title, prideDir.ConfigFile, prideDir.Path)
	if serr != nil {
		return serr
	}
	return nil
}

func OperationNewContent(destination string) *syserr.Err {
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	_, serr = site.CreateMarkdownFile(destination, site.GetDefaultMarkdownText(), dir.ConfigFile, dir.Path)
	if serr != nil {
		return serr
	}
	return nil
}
