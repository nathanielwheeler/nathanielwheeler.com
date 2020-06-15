package views

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	Yield interface{}
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
