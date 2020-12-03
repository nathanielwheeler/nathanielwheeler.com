package server

import (
	"testing"

	"github.com/matryer/is"
)

func TestParseTemplate(t *testing.T) {
	is := is.New(t)
  s := newServer()

  _, err := s.parseTemplates(nil, "pages/home")
  if err != nil {
    is.Fail()
  }
  is.NoErr(err)
}
