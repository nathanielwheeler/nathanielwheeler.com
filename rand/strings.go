package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 32

// RememberToken is a helper function designed to generate remember tokens of a set size
func RememberToken() (string, error) {
	return randString(rememberTokenBytes)
}

func randBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func randString(nBytes int) (string, error) {
	b, err := randBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
