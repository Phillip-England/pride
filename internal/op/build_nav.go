package op

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	fmt.Printf("🧬 building site navigation\n")
	config, serr := site.ConfigLoadFromCwd()
	if serr != nil {
		return serr
	}
	nav, serr := buildNav("./content", config.Theme, true)
	if serr != nil {
		return serr
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(nav))
	if err != nil {
		return syserr.New(syserr.CodeDev, fmt.Errorf("goquery failed to load generated navigation\nhere is their provided error: %s", err.Error()))
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
		innerNav := fmt.Sprintf("<nav nav-name='%s'>%s</nav>", navName, outerHtml)
		filePath := "./navigation/" + navName + ".html"
		err = os.WriteFile(filePath, []byte(innerNav), 0755)
		if err != nil {
			potErr = syserr.New(syserr.CodeDev, fmt.Errorf("failed to write generated nav to %s", filePath))
		}
		
	})
	if potErr != nil {
		return potErr
	}
	err = os.WriteFile("./navigation/default.html", []byte(nav), 0755)
	if err != nil {
		return syserr.New(syserr.CodeDev, fmt.Errorf("failed to write generated nav to %s", "./navigation/default.html"))
	}
	return nil
}

func buildNav(root string, theme string, isFirstPass bool) (string, syserr.SysErr) {
	var nav string
	base := filepath.Base(root)
	navName := strings.ReplaceAll(root, "/", "-")
	if strings.HasPrefix(navName, ".-") {
		navName = strings.Replace(navName, ".-", "", 1)
	}
	if isFirstPass {
		nav = fmt.Sprintf("<nav class='pride-nav' content-dir='%s' nav-name='default'><ul>", root)
	} else {
		nav = fmt.Sprintf(`<li class='pride-inner-nav' content-dir='%s' nav-name='%s'>%s<ul>`, root, base, navName)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", syserr.New(syserr.CodeDev, fmt.Errorf("failed to read %s while generating site navigation", root))
	}
	for _, entry := range entries {
		if entry.IsDir() {
			innerNav, serr := buildNav(filepath.Join(root, entry.Name()), theme, false)
			if serr != nil {
				return "", serr
			}
			nav += innerNav
		} else {
			mdFile, serr := site.MarkdownFileNew(filepath.Join(root, entry.Name()), "config", theme)
			if serr != nil {
				return "", serr
			}
			nav += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, mdFile.ServerPath, mdFile.Title)
		}
	}
	if isFirstPass {
		nav += "</nav></ul>"
	} else {
		nav += "</ul></li>"
	}
	return nav, nil
}