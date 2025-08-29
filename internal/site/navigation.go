package site

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Navigation struct {
	Menus map[string]Menu
}

func LoadNavigation(contentDir ContentDir) (Navigation, *syserr.Err) {
	var nav Navigation
	menus, serr := LoadMenus(contentDir)
	if serr != nil {
		return nav, serr
	}
	nav.Menus = menus
	fmt.Println(nav)
	return nav, nil
}

type Menu struct {
	Name  string
	Path  string
	Links []Link
}

func LoadMenus(contentDir ContentDir) (map[string]Menu, *syserr.Err) {
	menus := make(map[string]Menu)
	var potErr *syserr.Err
	potErr = nil
	err := filepath.WalkDir(contentDir.Path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		subDirName := strings.Replace(path, contentDir.Path, "", 1)
		name := ""
		if subDirName == "" {
			name = "Index"
		} else {
			if strings.Contains(subDirName, "/") && strings.Contains(subDirName, "\\") {
				potErr = syserr.New(syserr.Here(), "%s is an invalid content directory name, all content directory names must either contain \"/\" or \"/\" but never both", subDirName)
				return nil
			}
			var parts []string
			if strings.Contains(subDirName, "/") {
				parts = strings.Split(subDirName, "/")
			} else if strings.Contains(subDirName, "\\") {
				parts = strings.Split(subDirName, "\\")
			}
			if len(parts) > 1 {
				parts = parts[1:]
			}
			for _, part := range parts {
				runes := []rune(part)
				firstLetter := string(runes[0])
				var remainingLetters string
				if len(runes) > 1 {
					remainingLetters = string(runes[1:])
				}
				firstLetterStr := strings.ToUpper(firstLetter)
				name += firstLetterStr + remainingLetters
			}
		}
		var menu Menu
		menu.Name = name
		menu.Path = path
		links, serr := LoadLinks(menu.Path, contentDir)
		if serr != nil {
			potErr = serr
			return nil
		}
		menu.Links = links
		menus[name] = menu

		return nil
	})
	if err != nil {
		return menus, syserr.New(syserr.Here(), "%s", err.Error())
	}
	if potErr != nil {
		return menus, potErr
	}
	return menus, nil
}

type Link struct {
	Href string
	Html string
}

func LoadLinks(path string, contentDir ContentDir) ([]Link, *syserr.Err) {
	links := []Link{}
	for _, mdFile := range contentDir.MarkdownFiles {
		mdParentPath := filepath.Dir(mdFile.Path)
		if mdParentPath != path {
			continue
		}
		var link Link
		link.Href = mdFile.ServerPath
		link.Html = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, link.Href, mdFile.Title)
		links = append(links, link)
	}
	return links, nil
}
