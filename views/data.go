package views

import (
	"log"

	"nathanielwheeler.com/models"
)

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	User  *models.User
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
	AlertMsgGeneric = `Something went wrong, please try again.  If the problem persists, please contact me directly at <a href="mailto:contact@nathanielwheeler.com">contact@nathanielwheeler.com</a>.`
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
