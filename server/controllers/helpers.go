package controllers

import (
  "net/http"

  "github.com/gorilla/schema"
)

func parseForm(req *http.Request, dest interface{}) error {
  if err := req.ParseForm(); err != nil {
    return err
  }
  decoder := schema.NewDecoder()
  // IgnoreUnknownKeys so that we can use CSRF protection in our forms.
  decoder.IgnoreUnknownKeys(true)
  if err := decoder.Decode(dest, req.PostForm); err != nil {
    return err
  }

  return nil
}
