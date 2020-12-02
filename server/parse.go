package server

import (
	"fmt"
	"html/template"
	"path/filepath"
)

func (s *server) parseTemplates(files ...string) (tpl *template.Template, err error) {
	// Automatically adds app layout
	files = append(files, filepath.Join("layouts", "app"))
	for i, v := range files {
		file := files[i]
		file = filepath.Join("client", "templates", v)
		file = file + ".tpl"
		files[i] = file
	}
	s.logMsg(fmt.Sprintf("Files: %v\n", files))
	tpl, err = template.ParseFiles(files...)
	if err != nil {
		s.logErr("Error parsing template file", err)
		return nil, err
	}
	return tpl, nil
}

func (s *server) parseData(data interface{}) interface{} {
	switch data.(type) {
	case nil:
		return nil
	default:
		return data
	}
}
