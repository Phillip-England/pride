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

func (op *NewContent) Exec(c cmd.Cmd) *syserr.Err {
	cmdMake, ok := c.(*cmd.New)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	contentPath := "./content/" + cmdMake.ArgDestinationStripped
	fmt.Printf("🧬 making new content at %s.md\n", contentPath)
	_, serr := site.ConfigLoad()
	if serr != nil {
		return serr
	}
	contentDirPath := "./content"
	_, err := os.Stat(contentDirPath)
	if err != nil {
		return syserr.New(syserr.Here(), err.Error())
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
		return serr
	}
	return nil
}
