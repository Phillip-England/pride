package build

import (
	"fmt"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/syserr"
)

type HtmlFile struct {
	RootDir string
	Path string
}


func NewHtmlFile(rootDir string, route server.Route) (HtmlFile, *syserr.Err) {
	var f HtmlFile
	f.RootDir = rootDir
	f.Path = filepath.Join(rootDir, route.RelativePath)
	
	fmt.Println(f.Path)
	return f, nil
}