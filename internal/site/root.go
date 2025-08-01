package site

import (
	"os"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Root struct {
	PathRoot         string
	PathContent      string
	PathConfig       string
	PathStatic       string
	PathTemplates    string
	PathNavigation   string
	PathLayouts      string
	PathLayoutBase   string
	PathTemplateBase string
	PathCssBase      string
	PathContentBase  string
}

func RootNew(path string) Root {
	var dir Root
	dir.PathRoot = path
	dir.PathContent = path + "/content"
	dir.PathStatic = path + "/static"
	dir.PathTemplates = path + "/templates"
	dir.PathNavigation = path + "/navigation"
	dir.PathLayouts = path + "/layouts"
	dir.PathConfig = path + "/pride.toml"
	dir.PathLayoutBase = dir.PathLayouts + "/default.html"
	dir.PathTemplateBase = dir.PathTemplates + "/default.html"
	dir.PathCssBase = dir.PathStatic + "/default.css"
	dir.PathContentBase = dir.PathContent + "/index.md"
	return dir
}

func (dir Root) Create() *syserr.Err {
	err := os.Mkdir(dir.PathRoot, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return syserr.New(syserr.Here(), "<DESTINATION> %s already exists, please provide a new name or delete %s", dir.PathRoot, dir.PathRoot)
		}
		return syserr.New(syserr.Here(), "<DESTINATION> unanticipated error when creating %s, are you using a valid directory name for your site?", dir.PathRoot)

	}
	dirsToMake := []string{
		dir.PathContent,
		dir.PathStatic,
		dir.PathLayouts,
		dir.PathTemplates,
		dir.PathNavigation,
	}
	for _, dir := range dirsToMake {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return syserr.New(syserr.Here(), "unanticipated error when creating %s", dir)
		}
	}
	config := NewConfig(dir.PathConfig)
	serr := config.Create()
	if serr != nil {
		return serr
	}
	layout := LayoutNew(dir.PathLayoutBase)
	serr = layout.Create()
	if err != nil {
		return serr
	}
	template := TemplateNew(dir.PathTemplateBase)
	serr = template.Create()
	if serr != nil {
		return serr
	}
	css := CssNew(dir.PathCssBase)
	serr = css.Create()
	if serr != nil {
		return serr
	}
	md := ContentNew(dir.PathContentBase, "Home Page", false)
	serr = md.Create()
	if serr != nil {
		return serr
	}
	return nil
}
