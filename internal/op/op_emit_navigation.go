package op

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type OpEmitNavigation struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpEmitNavigation) Exec(c cmd.Cmd) *syserr.Err {
	_, ok := c.(*cmd.CmdEmit)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	targetDir, err := cmd.GetArg(3)
	if err != nil {
		return syserr.New(syserr.Here(), "missing <TARGETDIR> in `pride emit nav`")
	}
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	outDir, err := cmd.GetArg(4)
	if err != nil {
		return syserr.New(syserr.Here(), "missing <OUT> in `pride emit nav`")
	}
	_, serr := OperationEmitNavigation(targetDir, outDir)
	if serr != nil {
		return serr
	}
	return nil
}

func OperationEmitNavigation(targetDir string, out string) (site.NavigationMenu, *syserr.Err) {
	var navMenu site.NavigationMenu
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return navMenu, serr
	}
	for _, navMenu := range dir.Navigation.Menus {
		// success
		if navMenu.Path == targetDir {
			err := os.WriteFile(out, []byte(navMenu.Html), 0755)
			if err != nil {
				return navMenu, syserr.New(syserr.Here(), "%s", err.Error())
			}
			return navMenu, nil
		}
	}
	// err
	return navMenu, syserr.New(syserr.Here(), "failed to located a navigation menu for dir %s", targetDir)
}
