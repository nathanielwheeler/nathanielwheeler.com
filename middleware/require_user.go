package middleware

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
)

// RequireUser will hold UserService.
type RequireUser struct {
	models.UserService
}

// ApplyFn will return http.HandlerFunc that will check if a user is logged in then call next(res, req), or redirect them to the login page if they are not logged in.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")
		if err != nil {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}
		
		// Get context from request, make a new context from the existing one that has our user stored in it with the private user key, and create a new request from the existing one with the new context attached.
		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)

		next(res, req)
	})
}

// Apply will allow http.Handler interfaces to apply this middleware by passing ServeHTTP into ApplyFn
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
