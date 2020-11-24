package server

import (
	"log"
	"net/http"
	"time"

	"nathanielwheeler.com/server/services"
)

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	User  *services.User
	Yield interface{}
}

// PublicError is an interface applying to errors that have a Public method attached to them.
type PublicError interface {
	error
	Public() string
}

// Alert is used to render alert messages in templates
type Alert struct {
	Level   string
	Message string
}

const (
	// AlertLvlError indicates that an error has occurred.
	AlertLvlError = "danger"
	// AlertLvlWarning gives a warning to the user.
	AlertLvlWarning = "warning"
	// AlertLvlInfo indicates that some information has changed.
	AlertLvlInfo = "info"
	// AlertLvlSuccess indicates that an action was carried successfully.
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is the default message for unfiltered errors.
	AlertMsgGeneric = `Something went wrong, please try again.  If the problem persists, please contact me directly at <a href="mailto:nathan@mailftp.com">nathan@mailftp.com</a>.`
)

// SetAlert will create an alert using a constant error message
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok { // Type assertion, runs if error matches PublicError interface
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError takes in a message and constructs a custom error message
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// RedirectAlert redirects with a persisting error
func (d *Data) RedirectAlert(res http.ResponseWriter, req *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(res, alert)
	http.Redirect(res, req, urlStr, code)
}

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
		Level:   lvl.Value,
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
