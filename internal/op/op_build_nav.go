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

const defaultNavName = "default"
const navRootId = "pride-nav-root"
const navOuterClass = "pride-nav"
const navItemClass = "pride-nav-item"

type OpBuildNav struct {
	Code int
}

func (op *OpBuildNav) Exec(c cmd.Cmd) *syserr.Err {
	_, ok := c.(*cmd.CmdBuild)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure, did you use pointers correctly? did you return a valid op code?")
	}
	_, serr := BuildNavigation()
	if serr != nil {
		return serr
	}
	return nil
}

func BuildNavigation() (string, *syserr.Err) {
	prideDir, serr := site.LoadPrideDir()
	if serr != nil {
		return "", serr
	}
	nav, serr := build(prideDir.ConfigFile, true, prideDir.Path, prideDir.ContentDir.Path)
	if serr != nil {
		return "", serr
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(nav))
	if err != nil {
		return "", syserr.New(syserr.Here(), "%s", err.Error())
	}
	doc.Find(".pride-inner-nav").Each(func(i int, s *goquery.Selection) {
		if s.Find("a").Length() == 0 {
			s.Remove()
		}
	})
	var potErr *syserr.Err
	potErr = nil
	doc.Find(".pride-inner-nav").Each(func(i int, s *goquery.Selection) {
		inner := s.Children().First()
		outerHtml, err := goquery.OuterHtml(inner)
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
			return
		}
		navName, _ := s.Attr("nav-name")
		s = s.RemoveAttr("nav-name")
		s.RemoveClass("pride-inner-nav")
		innerNav := fmt.Sprintf("<nav class='pride-nav' nav-name='%s'>%s</nav>", navName, outerHtml)
		filePath := "./navigation" + "/" + navName + ".html"
		err = os.WriteFile(filePath, []byte(innerNav), 0755)
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
		}

	})
	if potErr != nil {
		return "", potErr
	}
	finalHtml, err := doc.Html()
	if err != nil {
		return "", syserr.New(syserr.Here(), "%s", err.Error())
	}
	err = os.WriteFile("./navigation"+"/default.html", []byte(finalHtml), 0755)
	if err != nil {
		return "", syserr.New(syserr.Here(), "%s", err.Error())
	}
	return nav, nil
}

func build(config site.ConfigFile, isFirstPass bool, prideDirPath string, currentDir string) (string, *syserr.Err) {
	var nav string
	var navName string
	pathDiff := strings.Replace(currentDir, prideDirPath, "", 1)

	// platform specific operation
	if strings.Contains(pathDiff, "/") {
		// linux/mac
		parts := strings.Split(pathDiff, "/")
		parts = parts[1:]
		if len(parts) == 1 {
			navName = "default"
		} else {
			navName = strings.Join(parts, "-")
		}
	} else {
		// windows
		parts := strings.Split(pathDiff, "\\")
		parts = parts[1:]
		if len(parts) == 1 {
			navName = "default"
		} else {
			navName = strings.Join(parts, "-")
		}
	}

	navName = strings.ReplaceAll(navName, "/", "-")
	navName = strings.Replace(navName, "content-", "", 1)
	if strings.HasPrefix(navName, ".-") {
		navName = strings.Replace(navName, ".-", "", 1)
	}
	navName = strings.TrimPrefix(navName, "-")
	if isFirstPass {
		nav = fmt.Sprintf("<nav id='%s' class='%s' nav-name='%s'><ul>", navRootId, navOuterClass, defaultNavName)
	} else {
		nav = fmt.Sprintf(`<li class='pride-inner-nav' nav-name='%s'>%s<ul>`, navName, navName)
	}
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return "", syserr.New(syserr.Here(), "%s", err.Error())
	}
	mdFiles := []site.MarkdownFile{}
	for _, entry := range entries {
		path := filepath.Join(currentDir, entry.Name())
		if entry.IsDir() {
			subDirPath := filepath.Join(currentDir, entry.Name())
			innerNav, err := build(config, false, prideDirPath, subDirPath)
			if err != nil {
				return "", err
			}
			nav += innerNav
		} else {
			mdFile, serr := site.LoadMarkdownFile(path, config.Theme, prideDirPath, currentDir)
			if serr != nil {
				return "", serr
			}
			if len(mdFiles) == 0 {
				if !mdFile.IsDraft {
					mdFiles = append(mdFiles, mdFile)
				}
				continue
			}
			prevMdFile := mdFiles[len(mdFiles)-1]
			prevDob, err := time.Parse(time.RFC3339, prevMdFile.Dob)
			if err != nil {
				return "", syserr.New(syserr.Here(), "%s", err.Error())
			}
			currentDob, err := time.Parse(time.RFC3339, mdFile.Dob)
			if err != nil {
				return "", syserr.New(syserr.Here(), "%s", err.Error())
			}
			if prevDob.Before(currentDob) {
				if !mdFile.IsDraft {
					mdFiles = append(mdFiles, mdFile)
				}
			} else {
				if !mdFile.IsDraft {
					mdFiles = append([]site.MarkdownFile{mdFile}, mdFiles...)
				}
			}
		}
	}
	for _, mdFile := range mdFiles {
		navItem := fmt.Sprintf(`<li class='%s' pride-dob='%s'><a href="%s">%s</a></li>`, navItemClass, config.Dob, mdFile.ServerPath, mdFile.Title)
		nav += navItem
	}
	if isFirstPass {
		nav += "</ul></nav>"
	} else {
		nav += "</ul></li>"
	}
	return nav, nil
}
