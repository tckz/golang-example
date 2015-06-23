package webapp

import (
	"./common"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type ParseRequest struct {
	layoutName   string
	templateName string
	chResult     chan<- ParseResult
}

type ParseResult struct {
	template *template.Template
	err      error
}

type TemplateParser struct {
	dirTemplate    string
	channelRequest chan<- ParseRequest
	watcher        *common.Watcher
}

func (self *TemplateParser) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"safeHtml": func(text string) template.HTML {
			return template.HTML(text)
		},
		"partial": func(name string, data interface{}) template.HTML {
			t, err := self.parse(name, "")
			if err != nil {
				panic(fmt.Errorf("*** Failed to parse: %v", err))
			}

			var b bytes.Buffer
			if err := t.ExecuteTemplate(&b, name, data); err != nil {
				fmt.Errorf("*** Failed to Execute: %v", err)
			}

			return template.HTML(b.Bytes())
		},
	}
}

func NewTemplateParser(dirTemplate string) *TemplateParser {
	parser := &TemplateParser{
		dirTemplate:   dirTemplate,
	}

	watcher, err := common.NewWatcher(dirTemplate)
	if err != nil {
		log.Fatalf("*** Failed to NewWatcher: %s", err)
	}
	parser.watcher = watcher

	channelNotify := make(chan bool)

	go func(chSendNotify chan<- bool) {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					break
				}
				log.Printf("Watch: %v, %s", ev.Op, ev.Name)
				chSendNotify <- true

			case err, ok := <-watcher.Errors:
				if !ok {
					break
				}
				if err != nil {
					log.Printf("*** watcher.Errors: %v", err)
				}
			}
		}
		close(chSendNotify)
	}(channelNotify)

	channelRequest := make(chan ParseRequest, 10)
	parser.channelRequest = channelRequest

	go func(chRequest <-chan ParseRequest, chNotify <-chan bool) {
		cache := make(map[string]ParseResult)

		for {
			select {

			case _, ok := <-chNotify:
				if !ok {
					break
				}
				cache = make(map[string]ParseResult)

			case req, ok := <-chRequest:
				if !ok {
					break
				}

				files := []string{}
				if req.layoutName != "" {
					files = append(files, filepath.Join(dirTemplate, req.layoutName))
				}
				if req.templateName != "" {
					files = append(files, filepath.Join(dirTemplate, req.templateName))
				}

				key := strings.Join(files, "\t")
				if found, ok := cache[key]; ok {
					req.chResult <- found
					continue
				}

				log.Printf("Parse: %s\n", key)
				t, err := template.New("").Funcs(parser.getFuncMap()).ParseFiles(files...)

				var result ParseResult
				if err != nil {
					result.err = err
				} else {
					for i, tt := range t.Templates() {
						//log.Printf("TT[%d]: %s, %#v", i, tt.Name(), tt.Tree)
						log.Printf("TT[%d]: %s", i, tt.Name())
					}
					result.template = t
				}
				req.chResult <- result
				cache[key] = result
			}
		}
	}(channelRequest, channelNotify)
	return parser
}

func (self *TemplateParser) parse(layout string, template string) (*template.Template, error) {
	ch := make(chan ParseResult)
	defer close(ch)

	var req = ParseRequest{
		layoutName:   layout,
		templateName: template,
		chResult:     ch,
	}
	self.channelRequest <- req

	ret, ok := <-ch
	if !ok {
		return nil, fmt.Errorf("*** No response: %s, %s", layout, template)
	} else if ret.err != nil {
		return nil, ret.err
	}

	return ret.template, nil
}

func (self *TemplateParser) Execute(layout string, template string, data interface{}) (*bytes.Buffer, error) {

	if t, err := self.parse(layout, template); err != nil {
		return nil, err
	} else {
		var b bytes.Buffer
		if err := t.ExecuteTemplate(&b, path.Base(layout), data); err != nil {
			return nil, fmt.Errorf("*** Failed to ExecuteTemplate: %v", err)
		}

		return &b, nil
	}
}

func (self *TemplateParser) Close() {
	close(self.channelRequest)
	self.watcher.Close()
}
