package handlers

import (
	"gotpl/lib"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// IHandleFunc same as calling mux.Handle but it instruments the endpoint so prometheus can expose the data
func IHandleFunc(mux *mux.Router, frag string, fn func(w http.ResponseWriter, req *http.Request)) {
	mux.Handle(frag, prometheus.InstrumentHandlerFunc(frag, fn))
}

// RenderFunc A higher order function that just renders a template with a specific name as a handlerfunc
func RenderFunc(tpl string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		Render(w, tpl, req)
	}
}

var tpl *template.Template
var funcs = template.FuncMap{}

// Render a template in the templates/ folder into a http ResponseWriter
func Render(w http.ResponseWriter, name string, data interface{}) {
	// Lets reparse the specific template when in dev mode so we dont have to restart the app
	if tpl == nil || !lib.CFG.Production {
		re := regexp.MustCompile(`templates[\\/]([A-Z\-_0-9a-z/\\]+).html`)
		tpl = template.New("tpl").Funcs(funcs)

		filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
			match := re.FindStringSubmatch(path)
			if len(match) > 0 {
				ar, err := ioutil.ReadFile(match[0])
				if err != nil {
					log.Panic(err)
				}

				_, err = tpl.New(strings.Replace(match[1], `\`, "/", 500)).Parse(string(ar))
				if err != nil {
					log.Panic(err)

				}
			}
			return nil
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf8")

	if err := tpl.ExecuteTemplate(w, name, data); err != nil {
		log.Panic(err)
	}
}
