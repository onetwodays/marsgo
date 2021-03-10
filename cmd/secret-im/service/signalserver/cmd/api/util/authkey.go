package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenAuthKey(salt string, password string) (authKey string) {
	h := sha1.New()
	h.Write([]byte(salt + password))
	authKey = hex.EncodeToString(h.Sum(nil))
	return
}