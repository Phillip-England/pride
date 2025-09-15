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
	Src site.PrideDir
	Dest string
}

func dirExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    return false
}


func GenerateBuild(dest string) (*syserr.Err) {
	// if the dest already exists, exit
	err := dirExists(dest)
	if err {
		return syserr.New(syserr.Here(), "<DESTINATION> %s already exists", dest)
	}
	// load up the prideDir
	prideDir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	// load up the server
	svr, serr := server.NewServer(8080, prideDir)
	if serr != nil {
		return serr
	}
	// go through each route, generating an html file for each one
	var htmlFiles []HtmlFile
	for _, route := range svr.Routes {
		htmlFile, serr := NewHtmlFile(dest, route, prideDir.ConfigFile, svr)
		if serr != nil {
			return serr
		}
		htmlFiles = append(htmlFiles, htmlFile)
	}
	// go through each html file and save it to disk
	for _, htmlFile := range htmlFiles {
		dir := filepath.Dir(htmlFile.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return syserr.New(syserr.Here(), "mkdir failed: %s", err.Error())
		}
		err := os.WriteFile(htmlFile.Path, []byte(htmlFile.Text), 0644)
		if err != nil {
			return syserr.New(syserr.Here(), "write failed: %s", err.Error())
		}
	}

	// copying over static assests 
	outStaticDir := filepath.Join(dest, "static")
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("application/xml", xml.Minify)
	var potErr *syserr.Err
	potErr = nil
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
		fmt.Println(outPath)
		// outPath := outStaticDir + "/" + path
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
		default:
			contentType = ""
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
		err = os.MkdirAll(filepath.Dir(outPath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}
		err = os.WriteFile(outPath, output, 0644)
		if err != nil {
			return fmt.Errorf("failed to write to %s: %w", outPath, err)
		}
		return nil
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

