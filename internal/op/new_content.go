package op

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type NewContent struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *NewContent) Exec(c cmd.Cmd) syserr.SysErr {
	cmdMake, ok := c.(*cmd.New)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure, did you use pointers correctly? did you return a valid op code?"))
	}
	contentPath := "./content/" + cmdMake.ArgDestinationStripped
	fmt.Printf("🧬 making new content at %s.md\n", contentPath)
	_, serr := site.ConfigLoadFromCwd()
	if serr != nil {
		return serr
	}
	contentDirPath := "./content"
	info, err := os.Stat(contentDirPath)
	if os.IsNotExist(err) || !info.IsDir() {
		return syserr.MiaNew(fmt.Errorf("could not located %s directory", contentDirPath))
	}
	if err != nil {
		return syserr.DevNew(fmt.Errorf("could not load information on %s after it has been confirmed to exist", contentDirPath))
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
	f := site.ContentNew(contentPath+".md", title, true)
	serr = f.Create()
	if serr != nil {
		if strings.Contains(serr.GetMessage(), "file exists") {
			return syserr.HelpNew(fmt.Errorf("%s already exists", contentPath))
		}
		return serr
	}
	return nil
}
