package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"nathanielwheeler.com/context"
)

/*
NOTE In this file, I am panicking so much because if these views are not parsed correctly, the entire app is useless.
*/

var (
	templateDir string = "views/"
	layoutDir   string = templateDir + "layouts/"
	templateExt string = ".html"
)

// View : Contains a pointer to a template and the name of a layout.
type View struct {
	Template *template.Template
	Layout   string
}

// NewView : Takes in a layout name, any number of filename strings, parses them into template, and returns the address of the new view.
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
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

// Render : Responsible for rendering the view.  Checks the underlying type of data passed into it.
func (v *View) Render(res http.ResponseWriter, req *http.Request, data interface{}) {
	res.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		// Done so I can access the data in a var with type Data.
		vd = d
	default:
		// If data is NOT of type Data, make one and set the data to the Yield field.
		vd = Data{
			Yield: data,
		}
	}
	// Lookup and set the user to the User field
	vd.User = context.User(req.Context())
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(res, `Something went wrong, please try again.  If the problem persists, please contact me directly at <a href="mailto:contact@nathanielwheeler.com">contact@nathanielwheeler.com</a>.`, http.StatusInternalServerError)
		return
	}
	io.Copy(res, &buf)
}

// ServeHTTP : Renders and serves views.
func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	v.Render(res, req, nil)
}

/*
	HELPERS
*/

// takes a slice of strings (should be file paths) and prepends the templateDir to each string
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// takes in a slice of strings (should be file paths) and appends the templateExt to each string
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
