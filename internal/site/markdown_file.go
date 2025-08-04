package site

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Phillip-England/pride/internal/syserr"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

func GetDefaultMarkdownText() string {
	return fmt.Sprintf(`+++
title = "Home Page"
dob = "%s"
draft = %t
template = "/templates/default.html"
+++

# A Header
Some Content
`, time.Now().UTC().Format(time.RFC3339), true)
}

type MarkdownFile struct {
	Path            string
	ServerPath      string
	Text            string
	Theme           string
	Html            string
	Meta            map[string]any
	FileName        string
	PathWithoutBase string
	Title           string
	Dob             string
	IsDraft         bool
	Template        string
}

func CreateMarkdownFile(path string, text string, configFile ConfigFile, prideDirPath string) (MarkdownFile, *syserr.Err) {
	var mdFile MarkdownFile
	file, err := os.Create(path)
	if err != nil {
		return mdFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	file.WriteString(text)
	loadedMdFile, serr := LoadMarkdownFile(path, configFile.Theme, prideDirPath)
	if serr != nil {
		return mdFile, serr
	}
	return loadedMdFile, nil
}

func LoadMarkdownFile(path string, theme string, prideRootDir string) (MarkdownFile, *syserr.Err) {
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
	template, ok := mdFile.Meta["template"].(string)
	if !ok {
		template = filepath.Join(prideRootDir, "templates", "default.html")
	} else {
		template = filepath.Join(prideRootDir, template)
	}
	_, err = os.Stat(template)
	if err != nil {
		return mdFile, syserr.New(syserr.Here(), "%s", err.Error())
	}
	mdFile.Template = template
	return mdFile, nil
}
