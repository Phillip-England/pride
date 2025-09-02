package site

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/PuerkitoBio/goquery"
)

// here is something we can easily do with navigation
// by default, all submenus are included in navigation
// this may not be what a user wants
// for example we may wish to use the praimry, root level navigation,
// without including all of the blog posts
// well, each navigation has a unique name based off it's path
// so, uses can target that name and then 'remove'
// an inner menu based off of it's name
// we just have to work out how to specify removal
// from the user's perspective
// accessing the submenus is easy due to the naming convention

type Navigation struct {
	Menus map[string]Menu
}

func dirNameToMenuName(contentDir ContentDir, path string) (string, *syserr.Err) {
	subDirName := strings.Replace(path, contentDir.Path, "", 1)
	if subDirName == "" {
		return "Index", nil
	}

	if strings.Contains(subDirName, "/") && strings.Contains(subDirName, "\\") {
		return "", syserr.New(
			syserr.Here(),
			"%s is an invalid content directory name, all content directory names must either contain \"/\" or \"\\\" but never both",
			subDirName,
		)
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

	name := ""
	for _, part := range parts {
		runes := []rune(part)
		firstLetter := string(runes[0])
		var remainingLetters string
		if len(runes) > 1 {
			remainingLetters = string(runes[1:])
		}
		name += strings.ToUpper(firstLetter) + remainingLetters
	}

	return name, nil
}

func LoadNavigation(contentDir ContentDir) (Navigation, *syserr.Err) {
	var nav Navigation
	nav.Menus = make(map[string]Menu)
	var potErr *syserr.Err

	filepath.Walk(contentDir.Path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		name, serr := dirNameToMenuName(contentDir, path)
		if serr != nil {
			potErr = serr
			return nil
		}

		menuHasMdFiles := false
		for _, mdFile := range contentDir.MarkdownFiles {
			mdFileParentPath := filepath.Dir(mdFile.Path)
			if mdFileParentPath == path {
				menuHasMdFiles = true
			}
		}

		if !menuHasMdFiles {
			return nil
		}

		menu, serr := LoadMenu(name, path, contentDir)
		if serr != nil {
			potErr = serr
			return nil
		}

		nav.Menus[menu.Name] = menu
		return nil
	})

	if potErr != nil {
		return nav, potErr
	}
	return nav, nil
}

type Menu struct {
	Name                string
	Path                string
	Html                string
	MarkdownParentDirs  []string
	MarkdownToParentMap map[string][]MarkdownFile
}

// 1. collect all the directories which have markdown files
// 2. filter markdown parent dirs outside the menu's scope
// 3. map each menu path to it's respective markdown files
// 4. generate navigation
// 5. generate primary nav (non-collapsable)
// 6. generating sub navigation
// 7. create a go query document and extract the nav and list items
// 8. go through each nav-item, and collect it's inner-text, then sort
// 9. go through each nav-item in sorted order, remove it, then re-add to end of list
// 10. retrieve the sorted html
// 11. add a <script> to enable sub-menu toggling on click

func LoadMenu(name string, menuPath string, contentDir ContentDir) (Menu, *syserr.Err) {
	var menu Menu
	menu.Name = name
	menu.Path = menuPath
	menu.MarkdownParentDirs = []string{}
	menu.MarkdownToParentMap = make(map[string][]MarkdownFile)

	// 1
	for _, mdFile := range contentDir.MarkdownFiles {
		mdParentDir := filepath.Dir(mdFile.Path)
		if mdParentDir == menu.Path {
			continue
		}
		if slices.Contains(menu.MarkdownParentDirs, mdParentDir) {
			continue
		}
		menu.MarkdownParentDirs = append(menu.MarkdownParentDirs, mdParentDir)
	}
	slices.Sort(menu.MarkdownParentDirs)
	menu.MarkdownParentDirs = append([]string{menu.Path}, menu.MarkdownParentDirs...)

	// 2
	filtered := []string{}
	for _, mdFileParentPath := range menu.MarkdownParentDirs {
		if !strings.Contains(mdFileParentPath, menu.Path) {
			continue
		}
		filtered = append(filtered, mdFileParentPath)
	}
	menu.MarkdownParentDirs = filtered

	// 3
	for _, navMenuPath := range menu.MarkdownParentDirs {
		for _, mdFile := range contentDir.MarkdownFiles {
			mdFileParent := filepath.Dir(mdFile.Path)
			if mdFileParent == navMenuPath {
				menu.MarkdownToParentMap[navMenuPath] = append(menu.MarkdownToParentMap[navMenuPath], mdFile)
			}
		}
	}

	// 4
	menu.Html = "<nav><ul>"

	for i, dir := range menu.MarkdownParentDirs {
		menuName, serr := dirNameToMenuName(contentDir, dir)
		if serr != nil {
			return menu, serr
		}

		mdFiles, ok := menu.MarkdownToParentMap[dir]
		if !ok {
			return menu, syserr.New(syserr.Here(), "key '%s' does not map to a []MarkdownFile", dir)
		}

		// 5
		if i == 0 {
			for _, mdFile := range mdFiles {
				menu.Html += "<li class='pride-nav-item'><a href=\"" + mdFile.ServerPath + "\">" + mdFile.Title + "</a></li>"
			}
			continue
		}

		// 6
		menu.Html += "<li class='pride-nav-item pride-nav-submenu' data-pride-submenu-name=" + menuName + ">"
		menu.Html += "<span class='pride-nav-submenu-title' style='cursor:pointer;'>" + menuName + "</span>"
		menu.Html += "<ul class='pride-nav-submenu' style='display:none;'>"
		for _, mdFile := range mdFiles {
			menu.Html += "<li><a href=\"" + mdFile.ServerPath + "\">" + mdFile.Title + "</a></li>"
		}
		menu.Html += "</ul>"
		menu.Html += "</li>"
	}

	menu.Html += "</ul></nav>"

	// 7
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(menu.Html))
	if err != nil {
		return menu, syserr.New(syserr.Here(), "%s", err.Error())
	}
	navMenu := doc.Find("body").Children().First()
	navMenuItems := navMenu.Find(".pride-nav-item")

	// 8
	foundTitles := []string{}
	navMenuItems.Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if slices.Contains(foundTitles, text) {
			return
		}
		foundTitles = append(foundTitles, text)
	})
	slices.Sort(foundTitles)

	// 9
	ul := navMenu.ChildrenFiltered("ul").First()
	for _, menuTitle := range foundTitles {
		navMenuItems.Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if menuTitle != text {
				return
			}
			itemClone := s.Clone()
			s.Remove()
			ul.AppendSelection(itemClone)
		})
	}

	// 10
	navMenuHtml, err := goquery.OuterHtml(navMenu)
	if err != nil {
		return menu, syserr.New(syserr.Here(), "%s", err.Error())
	}
	menu.Html = navMenuHtml

	// 11
	menu.Html += `
<script>
  (() => {
	let menus = document.querySelectorAll('.pride-nav-submenu')
	let titles = document.querySelectorAll('.pride-nav-submenu-title')
	for (let i = 0; i < titles.length; i++) {
	let title = titles[i]
	let menu  = menus[i]
	title.addEventListener('click', () => {
		let menuDisplay = menu.style.display
		if (!menuDisplay || menuDisplay === '' || menuDisplay === 'none') {
		// could replace this with a value provided in our config
		menu.style.display = 'block' 
		} else {
		menu.style.display = 'none'
		}
	})
	}
  })()
</script>

`

	return menu, nil
}

func (m *Menu) printMarkdownFiles() {
	for _, dir := range m.MarkdownParentDirs {
		fmt.Printf("Directory: %s\n", dir)

		files, ok := m.MarkdownToParentMap[dir]
		if !ok {
			fmt.Println("  (no files)")
			continue
		}

		for _, file := range files {
			fmt.Printf("  File: %s\n", file.Path)
		}
	}
}

func (m *Menu) Print() {
	fmt.Printf("Menu: %s\n", m.Name)
	fmt.Printf("Path: %s\n", m.Path)
	fmt.Println("Markdown Files:")
	m.printMarkdownFiles()
	fmt.Printf("Html: %s\n", m.Html)
}
