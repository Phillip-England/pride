package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type PrideDir struct {
	RootDir       string
	ContentDir    string
	ConfigPath    string
	StaticDir     string
	TemplatesDir  string
	NavigationDir string
	LayoutsDir    string
}

func NewPrideDir(path string) (PrideDir, *syserr.Err) {
	var dir PrideDir
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.RootDir = absolutePath
	dir.ContentDir = filepath.Join(absolutePath, "content")
	dir.ConfigPath = filepath.Join(absolutePath, "pride.toml")
	dir.StaticDir = filepath.Join(absolutePath, "static")
	dir.TemplatesDir = filepath.Join(absolutePath, "templates")
	dir.NavigationDir = filepath.Join(absolutePath, "navigation")
	dir.LayoutsDir = filepath.Join(absolutePath, "layouts")
	return dir, nil
}

func (dir PrideDir) CreateIfNotExists() *syserr.Err {
	dirsToMake := []string{
		dir.RootDir,
		dir.ContentDir,
		dir.StaticDir,
		dir.LayoutsDir,
		dir.TemplatesDir,
		dir.NavigationDir,
	}
	for _, d := range dirsToMake {
		err := os.MkdirAll(d, 0755)
		if err != nil {
			return syserr.New(syserr.Here(), "failed to create directory %s: %s", d, err.Error())
		}
	}
	config := NewConfig(dir.ConfigPath)
	serr := config.Create()
	if serr != nil {
		return serr
	}
	return nil
}

