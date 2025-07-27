package op

import (
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

	navDirs := []string{}
	code := syserr.ErrCode(0)
	err := filepath.Walk("./content", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			navDirs = append(navDirs, path)
		}
		return nil
	})
	if err != nil {
		return syserr.New(code, err)
	}

	navs := []string{}
	for _, dir := range navDirs {
		nav := fmt.Sprintf("<nav src='%s'><ul>", dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return syserr.New(syserr.CodeDev, fmt.Errorf("failed to read dir %s for nav generation", dir))
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			mdFile, err := site.MarkdownFileNew(path, "content", config.Theme)
			if err != nil {
				return err
			}
			nav += fmt.Sprintf("<li><a href='%s'>%s</a></li>", mdFile.ServerPath, mdFile.Title)
		}
		nav += "</ul></nav>"
		navs = append(navs, nav)
	}

	navNames := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for i, nav := range navs {
		os.WriteFile(navNames[i]+".html", []byte(nav), 0755)
	}

	return nil
}
