package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

/*
NOTE In this file, I am panicking so much because if these views are not parsed correctly, the entire app is useless.
TODO Make an error handling system
*/

// View : Contains a pointer to a template and the name of a layout.
type View struct {
	Template *template.Template
	Layout   string
}

// NewView : Takes in a layout name, any number of filename strings, parses them into template, and returns the address of the new view.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render : Responsible for rendering the view.
func (v *View) Render(res http.ResponseWriter, data interface{}) error {
	res.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(res, v.Layout, data)
}

// ServeHTTP : Renders and serves views.
func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if err := v.Render(res, nil); err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob("views/layouts/*.html")
	if err != nil {
		panic(err)
	}
	return files
}
