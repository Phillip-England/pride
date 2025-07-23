package pride

import (
	"fmt"
	"os"
	"strings"
)

type DirRoot struct {
	PathRoot         string
	PathContent      string
	PathConfig       string
	PathStatic       string
	PathTemplates    string
	PathLayouts      string
	PathLayoutBase   string
	PathTemplateBase string
	PathCssBase      string
	PathContentBase  string
}

func DirRootNew(path string) DirRoot {
	var dir DirRoot
	dir.PathRoot = path
	dir.PathContent = path + "/content"
	dir.PathStatic = path + "/static"
	dir.PathTemplates = path + "/templates"
	dir.PathLayouts = path + "/layouts"
	dir.PathConfig = path + "/pride.toml"
	dir.PathLayoutBase = dir.PathLayouts + "/default_layout.html"
	dir.PathTemplateBase = dir.PathTemplates + "/default_template.html"
	dir.PathCssBase = dir.PathStatic + "/default_styles.css"
	dir.PathContentBase = dir.PathContent + "/index.md"
	return dir
}

func (dir DirRoot) Create() SysErr {
	fmt.Printf("making new site in %s..\n", dir.PathRoot)
	err := os.Mkdir(dir.PathRoot, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> %s already exists, please provide a new name or delete %s", dir.PathRoot, dir.PathRoot))
		}
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("<DESTINATION> unanticipated error when creating %s, are you using a valid directory name for your site?", dir.PathRoot))

	}
	dirsToMake := []string{
		dir.PathContent,
		dir.PathStatic,
		dir.PathLayouts,
		dir.PathTemplates,
	}
	for _, dir := range dirsToMake {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when creating %s", dir))
		}
	}
	config := FileConfigNew(dir.PathConfig)
	syserr := config.Create()
	if syserr != nil {
		return syserr
	}
	layout := FileLayoutBaseNew(dir.PathLayoutBase)
	syserr = layout.Create()
	if syserr != nil {
		return syserr
	}
	template := FileTemplateBaseNew(dir.PathTemplateBase)
	syserr = template.Create()
	if syserr != nil {
		return syserr
	}
	css := FileCssBaseNew(dir.PathCssBase)
	syserr = css.Create()
	if syserr != nil {
		return syserr
	}
	md := FileContentBaseNew(dir.PathContentBase)
	syserr = md.Create()
	if syserr != nil {
		return syserr
	}
	return nil
}
