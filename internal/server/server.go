package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Server struct {
	LayoutsAndTemplates *template.Template
	Routes              []Route
	Port                int
	Mux                 *http.ServeMux
}

// 1. load all .html layouts found in the ./layouts dir
// 2. load all .html templates found in the ./templates dir
// 3. generate routes from .md files
// 4. ensure static files can be served from static
// 5. setup routes
// 6. serving the routes
func NewServer(port int, prideDir site.PrideDir) (Server, *syserr.Err) {
	var svr Server
	svr.Port = port
	svr.Mux = http.NewServeMux()
	tmpl := template.New("")
	// 1
	tmpl, serr := LoadLayouts(tmpl, prideDir.LayoutsDir.Path)
	if serr != nil {
		return svr, serr
	}
	// 2
	tmpl, serr = LoadTemplates(tmpl, prideDir.TemplatesDir.Path)
	if serr != nil {
		return svr, serr
	}
	for _, tmp := range tmpl.Templates() {
		fmt.Println(tmp.Name())
	}
	svr.LayoutsAndTemplates = tmpl
	svr.Routes = []Route{}
	for _, mdFile := range prideDir.ContentDir.MarkdownFiles {
		route, serr := NewRoute(mdFile)
		if serr != nil {
			return svr, serr
		}
		svr.Routes = append(svr.Routes, route)
	}
	// 3
	fs := http.FileServer(http.Dir(prideDir.StaticDir.Path))
	svr.Mux.Handle("GET /static", http.StripPrefix("/static/", fs))
	// 4
	for _, route := range svr.Routes {
		svr.Mux.HandleFunc("GET "+route.MarkdownFile.ServerPath, func(w http.ResponseWriter, r *http.Request) {
			err := svr.LayoutsAndTemplates.ExecuteTemplate(w, route.LayoutName, map[string]interface{}{
				"Meta":    route.MarkdownFile.Meta,
				"Content": route.MarkdownFile.Html,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	// 5
	portStr := strconv.Itoa(port)
	host := "localhost"
	addr := host + ":" + portStr
	fmt.Printf("serving on %s\n", addr)
	err := http.ListenAndServe(addr, svr.Mux)
	if err != nil {
		return svr, syserr.New(syserr.Here(), "%s", err.Error())
	}
	return svr, nil
}
