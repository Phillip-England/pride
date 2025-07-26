package op

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type BuildNav struct {
	Code int
}

func (op *BuildNav) Exec(c cmd.Cmd) syserr.SysErr {
	_, ok := c.(*cmd.Build)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure, did you use pointers correctly? did you return a valid op code?"))
	}
	fmt.Printf("🧬 building site navigation\n")
	config, serr := site.ConfigLoadFromCwd()
	if serr != nil {
		return serr
	}
	nav := "<nav>"
	modNav, _, serr := makeNav(nav, "./content", []string{}, config.Theme)
	if serr != nil {
		return serr
	}
	nav = modNav
	nav = nav + "</nav>"
	os.WriteFile("./nav.html", []byte(nav), 0755)
	return nil
}

func makeNav(nav string, rootPath string, handledPaths []string, theme string) (string, []string, syserr.SysErr) {
	var code syserr.ErrCode
	code = 0
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		for _, handledPath := range handledPaths {
			if path == handledPath {
				return nil
			}
		}
		handledPaths = append(handledPaths, path)
		if info.IsDir() {
			nav += "<ul>"
			modNav, modHandledPaths, serr := makeNav(nav, path, handledPaths, theme)
			if serr != nil {
				code = serr.GetCode()
				return errors.New(serr.GetMessage())
			}
			nav = modNav
			handledPaths = modHandledPaths
			nav += "</ul>"
		} else {
			mdFile, serr := site.MarkdownFileNew(path, "content", theme)
			if serr != nil {
				code = serr.GetCode()
				return errors.New(serr.GetMessage())
			}
			nav += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, mdFile.ServerPath, mdFile.Title)
		}
		return nil
	})
	if err != nil {
		return nav, handledPaths, syserr.New(syserr.ErrCode(code), err)
	}
	return nav, handledPaths, nil
}
