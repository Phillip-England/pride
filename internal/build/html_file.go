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
	RootDir string
	Path string
	Document *goquery.Document
	Text string
}


func NewHtmlFile(rootDir string, route server.Route, configFile site.ConfigFile) (HtmlFile, *syserr.Err) {
	var f HtmlFile
	f.RootDir = rootDir
	f.Path = filepath.Join(rootDir, route.RelativePath)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(route.MarkdownFile.Html))
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}
	f.Document = doc
	var potErr *syserr.Err
	// adjusting our hrefs to match the server we intend to deploy the site
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) == 0 {
			return
		}
		firstChar := string(href[0])
		if firstChar != "/" {
			return
		}
		if href == "/" {
			s.SetAttr("href", configFile.Server+"/index.html")
		} else {
			s.SetAttr("href", configFile.Server+href+".html")
		}
	})
	if potErr != nil {
		return f, potErr
	}
	// adjusting our links href's
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) == 0 {
			return
		}
		firstChar := string(href[0])
		if firstChar == "/" {
			s.SetAttr("href", configFile.Server+href)
		}
	})
	if potErr != nil {
		return f, potErr
	}
	// adjusting our script src's
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) == 0 {
			return
		}
		firstChar := string(src[0])
		if firstChar == "/" {
			s.SetAttr("src", configFile.Server+src)
		}
	})
	if potErr != nil {
		return f, potErr
	}
	// adjusting image tags
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) == 0 {
			return
		}
		firstChar := string(src[0])
		if firstChar == "/" {
			s.SetAttr("src", configFile.Server+src)
		}
	})
	if potErr != nil {
		return f, potErr
	}
	// adjusting iframe src's
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if len(src) == 0 {
			return
		}
		firstChar := string(src[0])
		if firstChar == "/" {
			s.SetAttr("src", configFile.Server+src)
		}
	})
	if potErr != nil {
		return f, potErr
	}
	// adjusting <object> data attr
	doc.Find("object").Each(func(i int, s *goquery.Selection) {
		data, _ := s.Attr("data")
		if len(data) == 0 {
			return
		}
		firstChar := string(data[0])
		if firstChar == "/" {
			s.SetAttr("data", configFile.Server+data)
		}
	})
	if potErr != nil {
		return f, potErr
	}
	htmlStr, err := doc.Html()
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	var buf bytes.Buffer
	err = m.Minify("text/html", &buf, bytes.NewBufferString(htmlStr))
	if err != nil {
		return f, syserr.New(syserr.Here(), "%s", err.Error())
	}
	f.Text = buf.String()
	return f, nil
}