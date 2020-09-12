package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"nathanielwheeler.com/context"

	"github.com/gorilla/csrf"
)

/*
NOTE In this file, I am panicking so much because if these views are not parsed correctly, the entire app is useless.
*/

var (
	templateDir  string = "views/"
	layoutDir    string = templateDir + "layouts/"
	componentDir string = templateDir + "components/"
	formDir      string = templateDir + "forms/"
	templateExt  string = ".html"
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
	files = append(files, componentFiles()...)
	files = append(files, formFiles()...)
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
  // Check cookie for alerts
  if alert := getAlert(req); alert != nil {
    vd.Alert = alert
    clearAlert(res)
  }
	// Lookup and set the user to the User field
	vd.User = context.User(req.Context())
	var buf bytes.Buffer

	// Create variables using current http request and add it onto the template FuncMap.
	csrfField := csrf.TemplateField(req)
	// pathPrefix := mux.Vars(req)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(res, `Something went wrong, please try again.  If the problem persists, please contact me directly at "contact@nathanielwheeler.com"`, http.StatusInternalServerError)
		return
	}
	io.Copy(res, &buf)
}

// ServeHTTP : Renders and serves views.
func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	v.Render(res, req, nil)
}

// RedirectAlert accepts params for an http.Redirect and performs a redirect after persisting the provided alert in a cookie.
func RedirectAlert(res http.ResponseWriter, req *http.Request, urlStr string, code int, alert Alert) {
  persistAlert(res, alert)
  http.Redirect(res, req, urlStr, code)
}

/*
  HELPERS
*/

// getAlert will retrieve an alert from cookie data.  Errors probably mean the alert is invalid, so will return nil.
func getAlert(req *http.Request) *Alert {
  lvl, err := req.Cookie("alert_level")
  if err != nil {
    return nil
  }
  msg, err := req.Cookie("alert_message")
  if err != nil {
    return nil
  }
  alert := Alert{
    Level: lvl.Value,
    Message: msg.Value,
  }
  return &alert
}

// persistAlert will use cookie data to store alerts, expiring after 5 minutes
func persistAlert(res http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
  }
  http.SetCookie(res, &lvl)
  http.SetCookie(res, &msg)
}

// clearAlert will clear alerts within cookies
func clearAlert(res http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
  }
  http.SetCookie(res, &lvl)
  http.SetCookie(res, &msg)
}

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

func componentFiles() []string {
	files, err := filepath.Glob(componentDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func formFiles() []string {
	files, err := filepath.Glob(formDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
