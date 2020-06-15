package data

import (
	"crypto/rand"
	"encoding/base64"
)

// Size of the token
const tokenBytes = 3

func getRandString(size int) (s string, err error) {
	b := make([]byte, size)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	s = base64.URLEncoding.EncodeToString(b)
	return
}

// GenToken returns a new token
func GenToken() (token string, err error) {
	token, err = getRandString(tokenBytes)
	return
}
