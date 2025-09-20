package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type PrideDir struct {
	Path         string
	ContentDir   ContentDir
	ConfigFile   ConfigFile
	StaticDir    StaticDir
	TemplatesDir TemplatesDir
	LayoutsDir   LayoutsDir
	Navigation   Navigation
}

func CreatePrideDir(path string) (PrideDir, error) {
	var dir PrideDir
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.Path = absolutePath

	if err := os.MkdirAll(dir.Path, 0755); err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}

	configFile, err := CreateConfigFile(filepath.Join(absolutePath, "pride.toml"))
	if err != nil {
		return dir, err
	}
	dir.ConfigFile = configFile

	contentDirPath := filepath.Join(absolutePath, "content")

	// NOTE: the order in which we create directories matters:
	// 1. static dir (no dependencies)
	// 2. templates dir (no dependencies)
	// 3. layouts dir (depends on templates)
	// 4. content dir (depends on templates)

	// 1
	staticDir, err := CreateStaticDir(filepath.Join(absolutePath, "static"))
	if err != nil {
		return dir, err
	}
	dir.StaticDir = staticDir

	// 2
	templatesDir, err := CreateTemplatesDir(filepath.Join(absolutePath, "templates"))
	if err != nil {
		return dir, err
	}
	dir.TemplatesDir = templatesDir

	// 3
	layoutsDir, err := CreateLayoutsDir(filepath.Join(absolutePath, "layouts"))
	if err != nil {
		return dir, err
	}
	dir.LayoutsDir = layoutsDir

	// 4
	contentDir, err := CreateContentDir(contentDirPath, configFile, dir.Path)
	if err != nil {
		return dir, err
	}
	dir.ContentDir = contentDir

	return dir, nil
}

func LoadPrideDir() (PrideDir, error) {
	var dir PrideDir

	config, err := LoadConfigFile()
	if err != nil {
		return dir, err
	}
	dir.ConfigFile = config

	absolutePath, err := filepath.Abs(filepath.Dir(config.Path))
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.Path = absolutePath

	contentDir, err := LoadContentDir(filepath.Join(dir.Path, "content"), config.Theme, absolutePath)
	if err != nil {
		return dir, err
	}
	dir.ContentDir = contentDir

	templatesDir := LoadTemplatesDir(filepath.Join(dir.Path, "templates"))
	dir.TemplatesDir = templatesDir

	layoutsDir, err := LoadLayoutsDir(filepath.Join(dir.Path, "layouts"))
	if err != nil {
		return dir, err
	}
	dir.LayoutsDir = layoutsDir

	staticDir, err := LoadStaticDir(filepath.Join(dir.Path, "static"))
	if err != nil {
		return dir, err
	}
	dir.StaticDir = staticDir

	nav, err := LoadNavigation(dir.ContentDir)
	if err != nil {
		return dir, err
	}
	dir.Navigation = nav

	return dir, nil
}
