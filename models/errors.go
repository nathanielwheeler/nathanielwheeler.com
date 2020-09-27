package models

import "strings"

const (
  errNotFound  modelError = "models: resource not found"
  
  errIDInvalid modelError = "models: ID provided was invalid"
  
	errUserIDRequired modelError = "models: user ID is required"

	errEmailRequired modelError = "models: email address is required"
	errEmailInvalid  modelError = "models: invalid email address"
	errEmailTaken    modelError = "models: email address is already taken"

	errPasswordRequired modelError = "models: password is required"
	errPasswordInvalid  modelError = "models: incorrect password"
	errPasswordTooShort modelError = "models: password was too short"
	// TODO password validator for max length (64)
	// TODO password validator for restricted characters in password

	errRememberRequired modelError = "models: remember token required"
  errRememberTooShort modelError = "models: remember token should be at least 32 bytes"

	errTitleRequired modelError = "models: title is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
