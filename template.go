package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var bufPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", errors.New("yield called inappropriately")
	},
}

var layout = template.Must(
	template.
		New("layout.html").
		Funcs(layoutFuncs).
		ParseFiles("templates/layout.html"),
)

var templates = template.Must(template.ParseGlob("templates/**/*.html"))

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

	data["CurrentUser"] = RequestUser(r)
	data["Flash"] = r.FormValue("flash")

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer bufPool.Put(buf)
			err := templates.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String()), err
		},
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	err := template.Must(layout.Clone()).Funcs(funcs).Execute(buf, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Aaaah! Something went wrong", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
