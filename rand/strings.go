package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 32

// Bytes will generate n random bytes.  Otherwise produces an error.  Uses the crypto/rand package, so perfect for remember tokens.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slice of size n and then return a base64 encoded string of that byte slice.
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function designed to generate remember tokens of a set byte size.
func RememberToken() (string, error) {
	return String(rememberTokenBytes)
}