package site

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/Phillip-England/pride/internal/syserr"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

var lowercaseWords = map[string]bool{
	"a": true, "an": true, "the": true,
	"and": true, "but": true, "or": true, "nor": true, "for": true, "so": true, "yet": true,
	"at": true, "by": true, "down": true, "from": true, "in": true,
	"into": true, "like": true, "near": true, "of": true, "off": true, "on": true,
	"onto": true, "out": true, "over": true, "past": true, "per": true, "plus": true,
	"than": true, "to": true, "up": true, "upon": true, "via": true, "with": true,
	"as": true, "if": true, "that": true, "when": true, "while": true,
}

func capitalize(word string) string {
	if word == "" {
		return ""
	}
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func TitleFromPath(path string) string {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	words := strings.Split(name, "-")
	for i, word := range words {
		lower := strings.ToLower(word)
		if i == 0 || i == len(words)-1 || !lowercaseWords[lower] {
			words[i] = capitalize(lower)
		} else {
			words[i] = lower
		}
	}
	return strings.Join(words, " ")
}

func GetDefaultMarkdownText() string {
	return fmt.Sprintf(`+++
title = "Home Page"
dob = "%s"
draft = %t
template = "/templates/default.html"
+++

# A Bold Title
Some lovely content.
`, time.Now().UTC().Format(time.RFC3339), true)
}

type MarkdownFile struct {
	Path            string
	ServerPath      string
	Text            string
	Theme           string
	Html            string
	Meta     map[string]any
	FileName        string
	PathWithoutBase string
	Title           string
	Dob             string
	IsDraft         bool
	LayoutPath      string
	Menus           []string
}

func CreateMarkdownFile(path string, title string, isDraft bool, menuNames []string, layoutPath string, configFile ConfigFile, prideDirPath string, contentDirPath string) (MarkdownFile, *syserr.Err) {
	var mdFile MarkdownFile
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return mdFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	file, err := os.Create(path)
	if err != nil {
		return mdFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	if title == "" {
		title = TitleFromPath(path)
	}
	menuNameStr := "["
	for i, menuName := range menuNames {
		menuNameStr += "\""
		menuNameStr += menuName
		menuNameStr += "\""
		if i == len(menuNames)-1 {
			continue
		}
		menuNameStr += ", "
	}
	menuNameStr += "]"
	markdownContent := fmt.Sprintf(`+++
title = "%s"
dob = "%s"
draft = %t
layout = "%s"
menus = %s
+++

# A Header
Some Content
` + "```go\nfmt.Println(\"Hello, World!\")\n```", title, time.Now().UTC().Format(time.RFC3339), isDraft, layoutPath, menuNameStr)
	file.WriteString(markdownContent)
	loadedMdFile, serr := LoadMarkdownFile(path, configFile.Theme, prideDirPath, contentDirPath)
	if serr != nil {
		return loadedMdFile, serr
	}
	return loadedMdFile, nil
}

func LoadMarkdownFile(path string, theme string, prideRootDir string, contentDirPath string) (MarkdownFile, *syserr.Err) {
	var mdFile MarkdownFile
	mdBytes, err := os.ReadFile(path)
	if err != nil {
		return mdFile, syserr.New(syserr.Here(), "failed to read %s", path)
	}
	mdFile.Text = string(mdBytes)
	mdFile.Path = path
	mdFile.Theme = theme
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&frontmatter.Extender{
				Mode: frontmatter.SetMetadata,
			},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(mdBytes, &buf, parser.WithContext(context)); err != nil {
		return mdFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	mdFile.Html = buf.String()
	mdFile.FileName = filepath.Base(mdFile.Path)
	mdFile.PathWithoutBase = strings.ReplaceAll(mdFile.Path, mdFile.FileName, "")
	root := md.Parser().Parse(text.NewReader(mdBytes))
	doc := root.OwnerDocument()
	mdFile.Meta = doc.Meta()
	title, ok := mdFile.Meta["title"].(string)
	if !ok {
		title = "Hello, World!"
	}
	mdFile.Title = title
	dob, ok := mdFile.Meta["dob"].(string)
	if !ok {
		dob = time.Now().UTC().Format(time.RFC3339)
	}
	mdFile.Dob = dob
	draft, ok := mdFile.Meta["draft"].(bool)
	if !ok {
		draft = true
	}
	mdFile.IsDraft = draft
	layout, ok := mdFile.Meta["layout"].(string)
	if !ok {
		layout = filepath.Join(prideRootDir, "layouts", "default.html")
	} else {
		layout = filepath.Join(prideRootDir, layout)
	}
	mdFile.LayoutPath = layout
	menus, ok := mdFile.Meta["menus"].([]string)
	if !ok {
		menus = []string{}
	}
	mdFile.Menus = menus
	// resolving the server path is platform specific
	trimmed := strings.ReplaceAll(path, prideRootDir, "")
	if strings.Contains(trimmed, "/") {
		// linux/macOS
		parts := strings.Split(trimmed, "/")
		parts = parts[2:]
		joined := strings.Join(parts, "/")
		joined = strings.TrimSuffix(joined, ".md")
		joined = "/" + joined
		if joined == "/index" {
			joined = "/"
		}
		mdFile.ServerPath = joined
	} else {
		// windows
		parts := strings.Split(trimmed, "\\")
		if len(parts) > 2 {
			parts = parts[2:]
		}
		joined := strings.Join(parts, "/")
		joined = strings.TrimSuffix(joined, ".md")
		joined = "/" + joined
		if joined == "/index" {
			joined = "/"
		}
		mdFile.ServerPath = joined
	}

	return mdFile, nil
}
