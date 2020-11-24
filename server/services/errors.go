package services

import "strings"

const (
  ErrNotFound  serviceError = "services: resource not found"
  
  errIDInvalid serviceError = "services: ID provided was invalid"
  
	errUserIDRequired serviceError = "services: user ID is required"

	errEmailRequired serviceError = "services: email address is required"
	errEmailInvalid  serviceError = "services: invalid email address"
	errEmailTaken    serviceError = "services: email address is already taken"

	ErrPasswordInvalid  serviceError = "services: incorrect password"
	errPasswordRequired serviceError = "services: password is required"
	errPasswordTooShort serviceError = "services: password was too short"
	// TODO password validator for max length (64)
	// TODO password validator for restricted characters in password

	errRememberRequired serviceError = "services: remember token required"
  errRememberTooShort serviceError = "services: remember token should be at least 32 bytes"

	errTitleRequired serviceError = "services: title is required"
)

type serviceError string

func (e serviceError) Error() string {
	return string(e)
}

func (e serviceError) Public() string {
	s := strings.Replace(string(e), "services: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
