package middleware

import (
	"net/http"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
)

// User middleware will lookup the current user via their remember token cookie using the UserService.  If found, they will be set on the request context.  Either way, the next handler is always called.
type User struct {
	models.UserService
}

// Apply will allow http.Handler interfaces to be handled by middleware by applying ServeHTTP to the handler and passing it into ApplyFn
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will take in an http.HandlerFunc and run middleware that will check for a remember token cookie and set it to the request context.  It will always call the next handler.
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")
		if err != nil {
			next(res, req)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(res, req)
			return
		}
		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)
		next(res, req)
	})
}

// RequireUser will redirect a user to /login if they are not logged in.  This middleware assumes that User middleware has already been run, otherwise it will always redirect users.
type RequireUser struct{}

// Apply will allow http.Handler interfaces to be handled by middleware by applying ServeHTTP to the handler and passing it into ApplyFn
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will take in an http.HandlerFunc and run middlware that will check for a remember token cookie.  If there is no user in context, the user will be redirected to /login.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		user := context.User(req.Context())
		if user == nil {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}
		next(res, req)
	})
}