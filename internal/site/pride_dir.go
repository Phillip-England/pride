package site

import (
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type PrideDir struct {
	Path       string
	ContentDir    ContentDir
	ConfigFile    ConfigFile
	StaticDir    StaticDir
	TemplatesDir  TemplatesDir
	NavigationDir NavigationDir
	LayoutsDir   LayoutsDir
}

func NewPrideDir(path string) (PrideDir, *syserr.Err) {
	var dir PrideDir
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.Path = absolutePath
	configFile, serr := NewConfigFile(filepath.Join(absolutePath, "pride.toml"))
	if serr != nil {
		return dir, serr
	}
	dir.ConfigFile = configFile
	contentDirPath := filepath.Join(absolutePath, "content")
	contentDir, serr := NewContentDir(contentDirPath)
	if serr != nil {
		return dir, serr
	}
	dir.ContentDir = contentDir
	templatesDir, serr := NewTemplatesDir(filepath.Join(absolutePath, "templates"))
	if serr != nil {
		return dir, serr
	}
	dir.TemplatesDir = templatesDir
	staticDir, serr := NewStaticDir(filepath.Join(absolutePath, "static"))
	if err != nil {
		return dir, serr
	}
	dir.StaticDir = staticDir
	navigationDir, serr := NewNavigationDir(filepath.Join(absolutePath, "navigation"))
	if err != nil {
		return dir, serr
	}
	dir.NavigationDir = navigationDir
	layoutsDir, serr := NewLayoutsDir(filepath.Join(absolutePath, "layouts"))
	if err != nil {
		return dir, serr
	}
	dir.LayoutsDir = layoutsDir
	return dir, nil
}


func LoadPrideDir() (PrideDir, *syserr.Err) {
	var dir PrideDir
	config, serr := LoadConfigFile()
	if serr != nil {
		return dir, serr
	}
	absolutePath, err := filepath.Abs(filepath.Dir(config.Path))
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.Path = absolutePath
	// need to actually load in each directory
	return dir, nil
}
