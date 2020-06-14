package data

import (
	"crypto/rand"
	"encoding/base64"
)

// Length of the random string in bytes and chars
// Due to base64: 4 * 20 / 3, rounded to multiple of 4
const tokenBytes = 6
const tokenChars = ((4 * tokenBytes / 3) + 3) & ^3

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
