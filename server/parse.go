package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
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

type markdown struct {
	Body     *bytes.Buffer
	MetaData map[string]interface{}
}

// parseMarkdown will take in a markdown file location and return html
func (s *server) parseMarkdown(file string) markdown {
	var (
		buf bytes.Buffer
		err error
	)
	file = file + ".md"

	src, err := ioutil.ReadFile(file)
	if err != nil {
		s.logErr("failed to read .md file", err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	ctx := parser.NewContext()
	err = md.Convert([]byte(src), &buf, parser.WithContext(ctx))
	if err != nil {
		s.logErr("markdown failed to parse", err)
	}

	return markdown{
		Body:     &buf,
		MetaData: meta.Get(ctx),
	}
}
