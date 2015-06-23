package main

import (
	"./webapp"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
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

type MyHTMLRender struct {
	layoutName   string
	templateName string
	data         interface{}
}

var htmlContentType = []string{"text/html; charset=utf-8"}

func (self MyHTMLRender) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)

	b, err := templateParser.Execute(self.layoutName, self.templateName, self.data)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	w.Write(b.Bytes())
	return nil
}

func renderTemplate(c *gin.Context, code int, layoutName string, templateName string, data interface{}) {
	r := MyHTMLRender{
		templateName: templateName,
		layoutName:   layoutName,
		data:         data,
	}
	c.Render(code, r)
}

func MyNotFound(c *gin.Context) {
	c.String(404, "Not foooound.")
}

func MyRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v\n%s", err, debug.Stack())
			renderTemplate(c, 500, "layout/default.tpl", "error.tpl", gin.H{})
		}
	}()
	c.Next()
}

func showIndex(c *gin.Context) {
	renderTemplate(c, 200, "layout/default.tpl", "index.tpl", 
		gin.H{
			"title":   c.Param("title"),
			"numbers": []int{1, 2, 3, 4},
		})
}

func main() {
	flag.Parse()

	cpus := runtime.NumCPU()
	fmt.Fprintf(os.Stderr, "cpus=%d\n", cpus)
	runtime.GOMAXPROCS(cpus)

	templateParser = webapp.NewTemplateParser("view")
	defer templateParser.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(MyRecover)
	router.Static("/assets", "./public/assets")
	router.GET("/oreore/:title", showIndex)
	router.GET("/oreore/", showIndex)
	router.GET("/oreore", showIndex)
	router.NoRoute(MyNotFound)
	router.Run("0.0.0.0:3000")
}
