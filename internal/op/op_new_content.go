package op

import (
	"os"
	"strings"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpNewContent struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpNewContent) Exec(c cmd.Cmd) *syserr.Err {
	cmdNew, ok := c.(*cmd.CmdNew)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	_, serr := OperationNewContent(cmdNew.ArgDestination, true, []string{})
	if serr != nil {
		return serr
	}
	return nil
}

// 1. new content must be unique
// 2. the names within menuNames must all begin with an uppercase letter
func OperationNewContent(destination string, isDraft bool, menuNames []string) (site.MarkdownFile, *syserr.Err) {
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return site.MarkdownFile{}, serr
	}
	// 1
	_, err := os.Stat(destination)
	if err == nil {
		return site.MarkdownFile{}, syserr.New(syserr.Here(), "%s already exists", destination)
	}
	// 2
	alteredMenuNames := []string{}
	for _, menuName := range menuNames {
		menuNameRunes := []rune(menuName)
		menuNameFirstChar := string(menuNameRunes[0])
		var menuNameRemainingRunes []rune
		if len(menuNameRunes) > 1 {
			menuNameRemainingRunes = menuNameRunes[1:]
		}
		menuNameFirstChar = strings.ToUpper(menuNameFirstChar)
		menuName = menuNameFirstChar + string(menuNameRemainingRunes)
		alteredMenuNames = append(alteredMenuNames, menuName)
	}
	mdFile, serr := site.CreateMarkdownFile(destination, "", isDraft, alteredMenuNames, "/layouts/default.html", dir.ConfigFile, dir.Path, dir.ContentDir.Path)
	if serr != nil {
		return mdFile, serr
	}
	return mdFile, nil
}
