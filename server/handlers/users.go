package handlers

import (
	"net/http"
	"time"

	"nathanielwheeler.com/util"
	views "nathanielwheeler.com/ui"
)

// NewUsers initializes the view for users
func NewUsers(us models.UserService) *Users {
	return &Users{
		RegisterView: views.NewView("app", "users/register"),
		LoginView:    views.NewView("app", "users/login"),
		us:           us,
	}
}

// Users holds reference for the Users view and service.
type Users struct {
	RegisterView *views.View
	LoginView    *views.View
	us           models.UserService
}

// Registration : GET /register
// — Renders a new registration form for a potential user
func (u *Users) Registration(res http.ResponseWriter, req *http.Request) {
	u.RegisterView.Render(res, req, nil)
}

// RegistrationForm is used to transform a webform into a registration request
type RegistrationForm struct {
	Email    string `schema:"email"`
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

// Register : POST /register
// — Used to process the signup form when a user tries to create a new user account
func (u *Users) Register(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form RegistrationForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		u.RegisterView.Render(res, req, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.RegisterView.Render(res, req, vd)
		return
	}
	err := u.signIn(res, &user)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// LoginForm is used to transform a webform into a login request
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login : POST /login
// — Used to process the login form when a user tries to log in as an existing user
func (u *Users) Login(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form LoginForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(res, req, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
		case models.ErrPasswordInvalid:
			vd.AlertError("Invalid email and/or password.")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(res, req, vd)
		return
	}

	err = u.signIn(res, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(res, req, vd)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// Logout : POST /logout
// — Used to process the logout form when a user chooses to logout.
func (u *Users) Logout(res http.ResponseWriter, req *http.Request) {
	// Expire user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	// Update user with a new remember token
	user := context.User(req.Context())
	token, _ := rand.RememberToken() // Ignoring errors because because unlikely and not much to do about it.
	user.Remember = token
	u.us.Update(user)

	http.Redirect(res, req, "/", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(res http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("remember_token")
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	_, err = u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	http.Redirect(res, req, "/", http.StatusFound)
}
