package server

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"nathanielwheeler.com/server/services"
	"nathanielwheeler.com/server/util"

	"github.com/gorilla/csrf"
)

var (
	templateDir string = "ui/templates/"
	templateExt string = ".tpl"
)

// UI : Contains a pointer to a template and the name of a layout.
type UI struct {
	Template *template.Template
	Layout   string
}

// NewUI : Takes in a layout name, any number of filename strings, parses them into template, and returns the address of the new view.
func NewUI(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, dirFiles("layouts")...)
	files = append(files, dirFiles("components")...)
	t, err := template.
		New("").
		Funcs(template.FuncMap{
			// csrfField is a placeholder for the gorilla/csrf field.  If it is not replaced in the render, it will throw an error.
			"csrfField": func() (template.HTML, error) {
				return "", errors.New("csrfField is not implemented")
			},
			// pathEscape will escape a path using the net/url package.
			"pathEscape": func(s string) string {
				return url.PathEscape(s)
			},
			// bodyFromPost will check for HTML in the post's body and, if it exists, render it.
			"bodyFromPost": func(post services.Post) template.HTML {
				if post.Body == "" {
					return template.HTML("Missing body...")
				}
				return template.HTML(post.Body)
			},
		}).
		ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render is responsible for rendering the view.  Checks the underlying type of data passed into it.  Then checks cookie for alerts, looks up the user, creates a CSRF field with the request data, and then executes the template.
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

	// Check cookie for alerts
	if alert := getAlert(req); alert != nil {
		vd.Alert = alert
		clearAlert(res)
	}

	// Lookup and set the user to the User field
	vd.User = util.User(req.Context())
	var buf bytes.Buffer

	// Create CSRF field using current http request and add it onto the template FuncMap.
	csrfField := csrf.TemplateField(req)

	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(res, `Something went wrong, please try again.  If the problem persists, please contact me directly at "nathan@mailftp.com"`, http.StatusInternalServerError)
		return
	}
	io.Copy(res, &buf)
}

// ServeHTTP renders and serves views.
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

func dirFiles(dir string) []string {
	files, err := filepath.Glob(templateDir + dir + "/*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
