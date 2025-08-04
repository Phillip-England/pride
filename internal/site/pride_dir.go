package site

import (
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/syserr"
)

type PrideDir struct {
	Path          string
	ContentDir    ContentDir
	ConfigFile    ConfigFile
	StaticDir     StaticDir
	TemplatesDir  TemplatesDir
	NavigationDir NavigationDir
	LayoutsDir    LayoutsDir
}

func CreatePrideDir(path string) (PrideDir, *syserr.Err) {
	var dir PrideDir
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return dir, syserr.New(syserr.Here(), "%s", err.Error())
	}
	dir.Path = absolutePath
	err = os.MkdirAll(dir.Path, 0755)
	if err != nil {
		syserr.New(syserr.Here(), "%s", err.Error())
	}
	configFile, serr := CreateConfigFile(filepath.Join(absolutePath, "pride.toml"))
	if serr != nil {
		return dir, serr
	}
	dir.ConfigFile = configFile
	contentDirPath := filepath.Join(absolutePath, "content")

	// NOTE: the order in which we create directories matters
	// this is because certain directories depend on each other
	// for example, all .md files in the content directory
	// will have an html template associated with them
	// we will then check the templates diretory to ensure the template
	// file actually exists, warning the user if they accidentally
	// choose an invalid template
	// creation should be done in the following order:
	// 1. static dir (no dependancies)
	// 2. templates dir (no dependancies)
	// 3. layouts dir (depends on templates)
	// 4. content dir (depends on templates)
	// 5. navigation dir (depends on content)

	// 1
	staticDir, serr := CreateStaticDir(filepath.Join(absolutePath, "static"))
	if serr != nil {
		return dir, serr
	}
	dir.StaticDir = staticDir

	// 2
	templatesDir, serr := CreateTemplatesDir(filepath.Join(absolutePath, "templates"))
	if serr != nil {
		return dir, serr
	}
	dir.TemplatesDir = templatesDir

	// 3
	layoutsDir, serr := CreateLayoutsDir(filepath.Join(absolutePath, "layouts"))
	if serr != nil {
		return dir, serr
	}
	dir.LayoutsDir = layoutsDir

	// 4
	contentDir, serr := CreateContentDir(contentDirPath, configFile, dir.Path)
	if serr != nil {
		return dir, serr
	}
	dir.ContentDir = contentDir

	// 5
	navigationDir, serr := CreateNavigationDir(filepath.Join(absolutePath, "navigation"))
	if serr != nil {
		return dir, serr
	}
	dir.NavigationDir = navigationDir

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
	contentDir, serr := LoadContentDir(filepath.Join(dir.Path, "content"))
	if serr != nil {
		return dir, serr
	}
	dir.ContentDir = contentDir
	navigationDir, serr := LoadNavigationDir(filepath.Join(dir.Path, "navigation"))
	if serr != nil {
		return dir, serr
	}
	dir.NavigationDir = navigationDir
	templatesDir, serr := LoadTemplatesDir(filepath.Join(dir.Path, "templates"))
	if serr != nil {
		return dir, serr
	}
	dir.TemplatesDir = templatesDir
	layoutsDir, serr := LoadLayoutsDir(filepath.Join(dir.Path, "layouts"))
	if serr != nil {
		return dir, serr
	}
	dir.LayoutsDir = layoutsDir
	staticDir, serr := LoadStaticDir(filepath.Join(dir.Path, "static"))
	if serr != nil {
		return dir, serr
	}
	dir.StaticDir = staticDir
	return dir, nil
}
