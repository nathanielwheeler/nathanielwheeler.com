package server

import (
	"html/template"
	"net/http"
	"sync"
)

// handleTemplate handles the execution of templates given data and any number of files.  The app layout is used as the template base.  A template from the pages directory must be used.  Components used must be included in this call.
func (s *server) handleTemplate(data interface{}, files ...string) http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = s.parseTemplates(w, files...)
			if err != nil {
				s.logErr("error parsing template files", err)
			}
		})
		w.Header().Set("Content-Type", "text/html")
		err = tpl.ExecuteTemplate(w, "app", s.parseData(nil))
		if err != nil {
			s.logErr("error executing template", err)
		}
	}
}
