package main

import (
	"./webapp"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

var templateParser *webapp.TemplateParser

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func MyRecover() martini.Handler {
	return func(c martini.Context, log *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%v\n%s", err, debug.Stack())

				val := c.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
				res := val.Interface().(http.ResponseWriter)

				res.WriteHeader(500)

				b, err := templateParser.Execute("layout/default.tpl", "error.tpl", nil)
				if err != nil {
					res.Write([]byte("Internal Server Error"))
				} else {
					res.Write(b.Bytes())
				}
			}
		}()

		c.Next()
	}
}

func MyNotFound() (int, string) {
	return 404, "Not foooound."
}

func showIndex(params martini.Params) (int, string) {
	b, err := templateParser.Execute("layout/default.tpl", "index.tpl", 
		map[string]interface{} {
			"title":   params["title"],
			"numbers": []int{1, 2, 3, 4},
		})
	if err != nil {
		panic(err)
	}

	return 200, string(b.Bytes())
}

func createMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Static("public"))
	m.Use(MyRecover())
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	return &martini.ClassicMartini{m, r}
}

func main() {
	flag.Parse()

	cpus := runtime.NumCPU()
	fmt.Fprintf(os.Stderr, "cpus=%d\n", cpus)
	runtime.GOMAXPROCS(cpus)

	templateParser = webapp.NewTemplateParser("view")
	defer templateParser.Close()

	m := createMartini()
	m.Get("/oreore/:title", showIndex)
	m.Get("/oreore/", showIndex)
	m.Get("/oreore", showIndex)
	m.NotFound(MyNotFound)
	m.RunOnAddr("0.0.0.0:3000")
}

