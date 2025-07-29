package site

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

type Content struct {
	Path    string
	Text    string
	Dob     string
	IsDraft bool
	Title   string
}

func ContentNew(path string, title string, isDraft bool) *Content {
	var f Content
	f.Title = title
	f.Path = path
	f.Dob = time.Now().UTC().Format(time.RFC3339)
	f.IsDraft = isDraft
	f.Text = fmt.Sprintf(`+++
title = "%s"
dob = "%s"
draft = %t
template = "/templates/default.html"
+++

# Welcome
This is the home page!
`, f.Title, f.Dob, f.IsDraft)
	return &f
}

func (f Content) Create() syserr.SysErr {
	dir := filepath.Dir(f.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("failed to create parent directories for %s => %s", f.Path, err.Error()))
	}
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when creating %s => %s", f.Path, err.Error()))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when writing to %s => %s", f.Path, err.Error()))
	}
	return nil
}

func GetContentPaths() ([]string, syserr.SysErr) {
	paths := []string{}
	potentialErrCode := 0
	err := filepath.WalkDir("./content", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return paths, syserr.New(syserr.ErrCode(potentialErrCode), err)
	}
	sort.Strings(paths)
	return paths, nil
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

func MarkdownFileLoad(path string, serverPrefix string, theme string, rootDir string) (*MarkdownFile, syserr.SysErr) {
	var mdFile MarkdownFile
	mdBytes, err := os.ReadFile(path)
	if err != nil {
		return &mdFile, syserr.New(syserr.CodeDev, fmt.Errorf("failed to read %s", path))
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
		return &mdFile, syserr.New(syserr.CodeLib, err)
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
		template = rootDir + "/templates/" + "default.html"
	} else {
		template = rootDir + template
	}
	_, err = os.Stat(template)
	if err != nil {
		return &mdFile, syserr.HelpNew(fmt.Errorf(`frontmatter 'template' value in %s pointed to a file which does not exist, please point the file towards a template found in ./templates`, path))
	}
	mdFile.Template = template
	serverPath := strings.TrimPrefix(mdFile.Path, serverPrefix)
	serverPath = strings.TrimSuffix(serverPath, ".md")
	if serverPath == "/index" {
		serverPath = "/"
	}
	if !strings.HasPrefix(serverPath, "/") {
		serverPath = "/" + serverPath
	}
	mdFile.ServerPath = serverPath
	return &mdFile, nil
}

func ContentLoadAll() ([]*MarkdownFile, syserr.SysErr) {
	content := []*MarkdownFile{}
	config, serr := ConfigLoadFromCwd()
	if serr != nil {
		return content, serr
	}
	paths, serr := GetContentPaths()
	if serr != nil {
		return content, serr
	}
	for _, path := range paths {
		mdFile, serr := MarkdownFileLoad(path, "content", config.Theme, config.TemplatesDir)
		if serr != nil {
			return content, serr
		}
		content = append(content, mdFile)
	}
	return content, nil
}

type ContentDir struct {
	Path    string
	MdFiles []MarkdownFile
}
