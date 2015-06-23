package main

import (
	"./webapp"
	"flag"
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

var templateParser *webapp.TemplateParser

func renderTemplate(w http.ResponseWriter, code int, layoutName string, templateName string, data interface{}) {

	b, err := templateParser.Execute(layoutName, templateName, data)
	if err != nil {
		panic(fmt.Sprintf("*** Failed to ExecuteTemplate: %v", err))
	}

	w.WriteHeader(code)
	w.Write(b.Bytes())
}

func showIndex(c web.C, w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, 200, "layout/default.tpl", "index.tpl",
		map[string]interface{}{
			"title":   c.URLParams["title"],
			"numbers": []int{1, 2, 3, 4},
		})
}

func MyRecoverer(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("%v\n%s", err, debug.Stack())
				renderTemplate(w, 500, "layout/default.tpl", "error.tpl", nil)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func MyNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
}

// $0 --bind 0.0.0.0:3000

func main() {
	flag.Parse()

	cpus := runtime.NumCPU()
	fmt.Fprintf(os.Stderr, "cpus=%d\n", cpus)
	runtime.GOMAXPROCS(cpus)

	templateParser = webapp.NewTemplateParser("view")
	defer templateParser.Close()

	goji.Get("/assets/*", http.FileServer(http.Dir("./public")))
	goji.Get("/oreore/:title", showIndex)
	goji.Get("/oreore/", showIndex)
	goji.Get("/oreore", showIndex)
	goji.NotFound(MyNotFound)
	goji.Use(MyRecoverer)

	graceful.PostHook(func() {
	})
	goji.Serve()
}
