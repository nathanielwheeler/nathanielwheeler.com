package views

import "html/template"

// NewView : Takes in any number of strings to create templates that have a standardized layout.
func NewView(layout string, files ...string) *View {
	files = append(files, "views/layouts/app.html")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View : Contains a pointer to a template.Template.
type View struct {
	Template *template.Template
	Layout   string
}
