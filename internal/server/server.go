package server

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Server struct {
	LayoutsAndTemplates *template.Template
	Routes              []*Route
	Port                int
	Mux                 *http.ServeMux
	Addr                string
	Html                string
}

// Middleware for logging requests in the format: [METHOD] [PATH] [DURATION]
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, elapsed)
	})
}

func NewServer(port int, prideDir site.PrideDir) (Server, error) {
	var svr Server
	svr.Port = port
	svr.Mux = http.NewServeMux()
	tmpl := template.New("")

	tmpl, err := LoadLayouts(tmpl, prideDir.LayoutsDir.Path)
	if err != nil {
		return svr, err
	}

	tmpl, err = LoadTemplates(tmpl, prideDir.TemplatesDir.Path)
	if err != nil {
		return svr, err
	}

	svr.LayoutsAndTemplates = tmpl
	svr.Routes = []*Route{}

	// --- Serve favicon at /favicon.ico ---
	faviconPath := prideDir.StaticDir.FaviconFile.Path
	svr.Mux.Handle("GET /favicon.ico", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})))

	// Static file server for everything else under /static/
	fs := http.FileServer(http.Dir(prideDir.StaticDir.Path))
	svr.Mux.Handle("GET /static/", loggingMiddleware(http.StripPrefix("/static/", fs)))

	// Load routes for markdown files
	for _, mdFile := range prideDir.ContentDir.MarkdownFiles {
		route, err := NewRoute(prideDir.Path, mdFile)
		if err != nil {
			return svr, err
		}
		svr.Routes = append(svr.Routes, route)
	}

	// Render templates for each route
	for _, route := range svr.Routes {
		var buf bytes.Buffer
		err := svr.LayoutsAndTemplates.ExecuteTemplate(&buf, route.LayoutName, map[string]interface{}{
			"Meta":    route.MarkdownFile.Meta,
			"Content": template.HTML(route.MarkdownFile.Html),
		})
		if err != nil {
			return svr, syserr.New(syserr.Here(), "%s", err.Error())
		}
		route.HtmlBytes = buf.Bytes()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if route.MarkdownFile.ServerPath == "/" && r.URL.Path != "/" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 Not Found"))
				return
			}
			w.Write(route.HtmlBytes)
		})

		svr.Mux.Handle("GET "+route.MarkdownFile.ServerPath, loggingMiddleware(handler))
	}

	portStr := strconv.Itoa(port)
	svr.Addr = "localhost:" + portStr
	return svr, nil
}
