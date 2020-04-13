package views

import (
	"net/http"
	"html/template"
	"path/filepath"
)

/*
NOTE In this file, I am panicking so much because if these views are not parsed correctly, the entire app is useless.
*/

// View : Contains a pointer to a template and the name of a layout.
type View struct {
	Template *template.Template
	Layout   string
}

// Render : Responsible for rendering the view.
func (v *View) Render(res http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(res, v.Layout, data)
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

func layoutFiles() []string {
	files, err := filepath.Glob("views/layouts/*.html")
	if err != nil {
		panic(err)
	}
	return files
}
