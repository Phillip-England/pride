package site

import (
	"fmt"
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

func (dir Root) Create() syserr.SysErr {
	err := os.Mkdir(dir.PathRoot, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return syserr.New(syserr.CodeHelp, fmt.Errorf("<DESTINATION> %s already exists, please provide a new name or delete %s", dir.PathRoot, dir.PathRoot))
		}
		return syserr.New(syserr.CodeHelp, fmt.Errorf("<DESTINATION> unanticipated error when creating %s, are you using a valid directory name for your site?", dir.PathRoot))

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
			return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when creating %s", dir))
		}
	}
	config := ConfigNew(dir.PathConfig)
	syserr := config.Create()
	if syserr != nil {
		return syserr
	}
	layout := LayoutNew(dir.PathLayoutBase)
	syserr = layout.Create()
	if syserr != nil {
		return syserr
	}
	template := TemplateNew(dir.PathTemplateBase)
	syserr = template.Create()
	if syserr != nil {
		return syserr
	}
	css := CssNew(dir.PathCssBase)
	syserr = css.Create()
	if syserr != nil {
		return syserr
	}
	md := ContentNew(dir.PathContentBase, "Home Page", false)
	syserr = md.Create()
	if syserr != nil {
		return syserr
	}
	return nil
}
