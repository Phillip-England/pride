package op

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/PuerkitoBio/goquery"
)

type BuildNav struct {
	Code int
}

func (op *BuildNav) Exec(c cmd.Cmd) syserr.SysErr {
	_, ok := c.(*cmd.Build)
	if !ok {
		return syserr.New(syserr.CodeDev, fmt.Errorf("type assertion failure, did you use pointers correctly? did you return a valid op code?"))
	}
	_, serr := BuildNavigation()
	if serr != nil {
		return serr
	}
	return nil
}

func BuildNavigation() (string, syserr.SysErr) {
	fmt.Printf("🧬 building site navigation\n")
	config, serr := site.ConfigLoadFromCwd()
	if serr != nil {
		return "", serr
	}
	nav, serr := build("./content", config, true)
	if serr != nil {
		return "", serr
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(nav))
	if err != nil {
		return "", syserr.New(syserr.CodeLib, fmt.Errorf("goquery failed to load generated navigation: %s", err.Error()))
	}
	var potErr syserr.SysErr
	potErr = nil
	doc.Find(".pride-inner-nav").Each(func(i int, s *goquery.Selection) {
		inner := s.Children().First()
		outerHtml, err := goquery.OuterHtml(inner)
		if err != nil {
			potErr =  syserr.New(syserr.CodeDev, fmt.Errorf("goquery failed to parse the outerHTML of an inner navigation, here is their provided error: %s", err.Error()))
			return
		}
		navName, _ := s.Attr("nav-name")
		s = s.RemoveAttr("nav-name")
		s = s.RemoveClass("pride-inner-nav")
		innerNav := fmt.Sprintf("<nav class='pride-nav' nav-name='%s'>%s</nav>", navName, outerHtml)
		filePath := "./navigation/" + navName + ".html"
		err = os.WriteFile(filePath, []byte(innerNav), 0755)
		if err != nil {
			potErr = syserr.New(syserr.CodeDev, fmt.Errorf("failed to write generated nav to %s because:\n%s", filePath, err.Error()))
		}
		
	})
	if potErr != nil {
		return "", potErr
	}
	err = os.WriteFile("./navigation/default.html", []byte(nav), 0755)
	if err != nil {
		return "", syserr.New(syserr.CodeDev, fmt.Errorf("failed to write generated nav to %s because:\n%s", "./navigation/default.html", err.Error()))
	}
	return nav, nil
}


func build(root string, config site.Config, isFirstPass bool) (string, syserr.SysErr) {
	var nav string
	navName := strings.ReplaceAll(root, "/", "-")
	navName = strings.Replace(navName, "content-", "", 1)
	if strings.HasPrefix(navName, ".-") {
		navName = strings.Replace(navName, ".-", "", 1)
	}
	if isFirstPass {
		nav = "<nav id='pride-nav-root' class='pride-nav' nav-name='default'><ul>"
	} else {
		nav = fmt.Sprintf(`<li class='pride-inner-nav' nav-name='%s'>%s<ul>`, navName, navName)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", syserr.New(syserr.CodeDev, fmt.Errorf("failed to read %s while generating site navigation", root))
	}
	mdFiles := []site.MarkdownFile{}
	for _, entry := range entries {
		if entry.IsDir() {
			innerNav, serr := build(filepath.Join(root, entry.Name()), config, false)
			if serr != nil {
				return "", serr
			}
			nav += innerNav
		} else {
			mdFile, serr := site.MarkdownFileNew(filepath.Join(root, entry.Name()), "content", config.Theme)
			if serr != nil {
				return "", serr
			}
			// navItem := fmt.Sprintf(`<li class='pride-nav-item' pride-dob='%s'><a href="%s">%s</a></li>`, config.Dob, mdFile.ServerPath, mdFile.Title)
			// nav += navItem
			if len(mdFiles) == 0 {
				mdFiles = append(mdFiles, mdFile)
				continue
			}
			prevMdFile := mdFiles[len(mdFiles)-1]
			prevDob, err := time.Parse(time.RFC3339, prevMdFile.Dob)
			if err != nil {
				return "", syserr.DevNew(fmt.Errorf("failed to parse time from nav link during nav generation\n%s", err.Error()))
			}
			currentDob, err := time.Parse(time.RFC3339, mdFile.Dob)
			if err != nil {
				return "", syserr.DevNew(fmt.Errorf("failed to parse time from nav link during nav generation\n%s", err.Error()))
			}
			if prevDob.Before(currentDob) {
				mdFiles = append(mdFiles, mdFile)
			} else {
				mdFiles = append([]site.MarkdownFile{mdFile}, mdFiles...)
			}
 		}
	}
	for _, mdFile := range mdFiles {
		navItem := fmt.Sprintf(`<li class='pride-nav-item' pride-dob='%s'><a href="%s">%s</a></li>`, config.Dob, mdFile.ServerPath, mdFile.Title)
		nav += navItem
	}
	if isFirstPass {
		nav += "</nav></ul>"
	} else {
		nav += "</ul></li>"
	}
	return nav, nil
}