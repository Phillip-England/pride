package build

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type HtmlFile struct {
	RootDir  string
	Path     string
	Document *goquery.Document
	Text     string
}

func chopFirstDir(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) <= 1 {
		return path
	}
	return strings.Join(parts[1:], "/")
}

func NewHtmlFile(rootDir string, route *server.Route, configFile site.ConfigFile, svr server.Server) (HtmlFile, error) {
	var f HtmlFile
	f.RootDir = rootDir
	f.Path = filepath.Join(rootDir, strings.TrimSuffix(chopFirstDir(route.RelativePath), ".md")+".html")

	routeHtml := string(route.HtmlBytes)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(routeHtml))
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}
	f.Document = doc

	// Adjust all resource URLs to point to the deployment server
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) == 0 || href[0] != '/' {
			return
		}
		if href == "/" {
			s.SetAttr("href", configFile.Server+"/index.html")
		} else {
			s.SetAttr("href", configFile.Server+href+".html")
		}
	})

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) != 0 && href[0] == '/' {
			s.SetAttr("href", configFile.Server+href)
		}
	})

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) != 0 && src[0] == '/' {
			s.SetAttr("src", configFile.Server+src)
		}
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) != 0 && src[0] == '/' {
			s.SetAttr("src", configFile.Server+src)
		}
	})

	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) != 0 && src[0] == '/' {
			s.SetAttr("src", configFile.Server+src)
		}
	})

	doc.Find("object").Each(func(i int, s *goquery.Selection) {
		data, _ := s.Attr("data")
		if len(data) != 0 && data[0] == '/' {
			s.SetAttr("data", configFile.Server+data)
		}
	})

	// Minify the final HTML
	htmlStr, err := doc.Html()
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	var buf bytes.Buffer
	if err := m.Minify("text/html", &buf, bytes.NewBufferString(htmlStr)); err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}

	f.Text = buf.String()
	return f, nil
}
