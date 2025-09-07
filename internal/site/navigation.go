package site

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

type Navigation struct {
	Menus map[string]NavigationMenu
}

func DirNameToMenuName(contentDir ContentDir, path string) (string, *syserr.Err) {
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
	nav.Menus = make(map[string]NavigationMenu)
	var potErr *syserr.Err

	filepath.Walk(contentDir.Path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		name, serr := DirNameToMenuName(contentDir, path)
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

type NavigationMenu struct {
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
// 12. wrap the nav in a valid go html {{ define }} clause
func LoadMenu(name string, menuPath string, contentDir ContentDir) (NavigationMenu, *syserr.Err) {
	var menu NavigationMenu
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
	menu.Html = "<nav class='pride-nav'><ul>"

	for i, dir := range menu.MarkdownParentDirs {
		menuName, serr := DirNameToMenuName(contentDir, dir)
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
		menu.Html += "<li class='pride-nav-item' data-pride-submenu-name=" + menuName + ">"
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
	navMenuHtml, err := PrettyFromSelection(navMenu)
	if err != nil {
		return menu, syserr.New(syserr.Here(), "%s", err.Error())
	}
	menu.Html = navMenuHtml

	// 11
	menu.Html += `
	<script>
	  (() => {
	   	let navMenus = document.querySelectorAll('.pride-nav')
		for (let i = 0; i < navMenus.length; i++) {
			let navMenu = navMenus[i]
			if (navMenu.getAttribute('pride-complete') == 'true') {
				continue
			}
			navMenu.setAttribute('pride-complete', 'true')
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
		}

	  })()
	</script>

	`

	// 12
	menu.Html = "{{ define \"" + "Nav" + menu.Name + "\"}}" + menu.Html + "{{ end }}"

	return menu, nil
}

func (m *NavigationMenu) printMarkdownFiles() {
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

func (m *NavigationMenu) Print() {
	fmt.Printf("Menu: %s\n", m.Name)
	fmt.Printf("Path: %s\n", m.Path)
	fmt.Println("Markdown Files:")
	m.printMarkdownFiles()
	fmt.Printf("Html: %s\n", m.Html)
}

// NormalizeHTML parses, normalizes, and pretty-prints a fragment (like your <nav>…</nav>)
func NormalizeHTML(input string) (string, error) {
	ctx := &html.Node{Type: html.ElementNode, Data: "div"}
	nodes, err := html.ParseFragment(strings.NewReader(input), ctx)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for i, n := range nodes {
		if err := renderPretty(&buf, n, 0); err != nil {
			return "", err
		}
		if i < len(nodes)-1 {
			buf.WriteString("\n")
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

func renderPretty(w io.Writer, n *html.Node, depth int) error {
	indent := strings.Repeat("  ", depth)

	switch n.Type {
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderPretty(w, c, depth); err != nil {
				return err
			}
		}
	case html.ElementNode:
		fmt.Fprint(w, indent)
		fmt.Fprint(w, "<", n.Data)
		for _, a := range n.Attr {
			fmt.Fprintf(w, ` %s="%s"`, a.Key, html.EscapeString(a.Val))
		}
		fmt.Fprint(w, ">")

		if hasNonEmptyChild(n) {
			fmt.Fprint(w, "\n")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := renderPretty(w, c, depth+1); err != nil {
					return err
				}
			}
			fmt.Fprintf(w, "%s</%s>\n", indent, n.Data)
		} else {
			fmt.Fprintf(w, "</%s>\n", n.Data)
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Fprintf(w, "%s%s\n", indent, text)
		}
	case html.CommentNode:
		fmt.Fprintf(w, "%s<!--%s-->\n", indent, strings.TrimSpace(n.Data))
	default:
	}
	return nil
}

func hasNonEmptyChild(n *html.Node) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			if strings.TrimSpace(c.Data) != "" {
				return true
			}
		case html.ElementNode, html.CommentNode:
			return true
		}
	}
	return false
}

func PrettyFromSelection(sel *goquery.Selection) (string, error) {
	if sel == nil || sel.Length() == 0 {
		return "", nil
	}
	var raw strings.Builder
	for _, n := range sel.Nodes {
		var buf bytes.Buffer
		if err := html.Render(&buf, n); err != nil {
			return "", err
		}
		raw.WriteString(buf.String())
	}
	return gohtml.Format(raw.String()), nil
}
