package server

import (
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
	Addr string
}


func NewServer(port int, prideDir site.PrideDir) (Server, *syserr.Err) {
	var svr Server
	svr.Port = port
	svr.Mux = http.NewServeMux()
	tmpl := template.New("")
	//
	tmpl, serr := LoadLayouts(tmpl, prideDir.LayoutsDir.Path)
	if serr != nil {
		return svr, serr
	}
	//
	tmpl, serr = LoadTemplates(tmpl, prideDir.TemplatesDir.Path)
	if serr != nil {
		return svr, serr
	}
	svr.LayoutsAndTemplates = tmpl
	svr.Routes = []Route{}
	for _, mdFile := range prideDir.ContentDir.MarkdownFiles {
		route, serr := NewRoute(prideDir.Path, mdFile)
		if serr != nil {
			return svr, serr
		}
		svr.Routes = append(svr.Routes, route)
	}
	//
	fs := http.FileServer(http.Dir(prideDir.StaticDir.Path))
	svr.Mux.Handle("GET /static", http.StripPrefix("/static/", fs))
	//
	for _, route := range svr.Routes {
		svr.Mux.HandleFunc("GET "+route.MarkdownFile.ServerPath, func(w http.ResponseWriter, r *http.Request) {
			err := svr.LayoutsAndTemplates.ExecuteTemplate(w, route.LayoutName, map[string]interface{}{
				"Meta":    route.MarkdownFile.Meta,
				"Content": template.HTML(route.MarkdownFile.Html),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	//
	portStr := strconv.Itoa(port)
	host := "localhost"
	addr := host + ":" + portStr
	svr.Addr = addr
	return svr, nil
}
