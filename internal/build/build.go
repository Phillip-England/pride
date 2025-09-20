package build

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

type Build struct {
	Src  site.PrideDir
	Dest string
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func GenerateBuild(dest string) error {
	// if the dest already exists, exit
	if dirExists(dest) {
		return syserr.New(syserr.Here(), "<DESTINATION> %s already exists", dest)
	}

	// load up the prideDir
	prideDir, err := site.LoadPrideDir()
	if err != nil {
		return err
	}

	// load up the server
	svr, err := server.NewServer(8080, prideDir)
	if err != nil {
		return err
	}

	// go through each route, generating an html file for each one
	var htmlFiles []HtmlFile
	for _, route := range svr.Routes {
		if route.MarkdownFile.IsDraft {
			continue
		}
		htmlFile, err := NewHtmlFile(dest, route, prideDir.ConfigFile, svr)
		if err != nil {
			return err
		}
		htmlFiles = append(htmlFiles, htmlFile)
	}

	// go through each html file and save it to disk
	for _, htmlFile := range htmlFiles {
		dir := filepath.Dir(htmlFile.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return syserr.New(syserr.Here(), "mkdir failed: %s", err.Error())
		}
		if err := os.WriteFile(htmlFile.Path, []byte(htmlFile.Text), 0644); err != nil {
			return syserr.New(syserr.Here(), "write failed: %s", err.Error())
		}
	}

	// copying over static assets
	outStaticDir := filepath.Join(dest, "static")
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("application/xml", xml.Minify)

	var potErr error
	filepath.Walk(prideDir.StaticDir.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(prideDir.StaticDir.Path, path)
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
			return nil
		}

		outPath := filepath.Join(outStaticDir, relativePath)
		input, err := os.ReadFile(path)
		if err != nil {
			potErr = syserr.New(syserr.Here(), "%s", err.Error())
			return nil
		}

		ext := filepath.Ext(path)
		var contentType string
		switch ext {
		case ".css":
			contentType = "text/css"
		case ".js":
			contentType = "text/javascript"
		case ".svg":
			contentType = "image/svg+xml"
		case ".html", ".htm":
			contentType = "text/html"
		case ".xml":
			contentType = "application/xml"
		}

		var output []byte
		if contentType != "" {
			output, err = m.Bytes(contentType, input)
			if err != nil {
				return fmt.Errorf("failed to minify %s: %w", path, err)
			}
		} else {
			output = input
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}
		if err := os.WriteFile(outPath, output, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %w", outPath, err)
		}

		return nil
	})

	if potErr != nil {
		return potErr
	}

	return nil
}
