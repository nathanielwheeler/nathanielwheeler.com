package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestHandleTemplate(t *testing.T) {
	// NOTE Frustratingly, tests treat the package as the working directory.  Therefore, I need to change to the real wd.
  os.Chdir("../") // Seriously, who thought that was a good idea?!
  
	is := is.New(t)
	s := newServer()


	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
}
