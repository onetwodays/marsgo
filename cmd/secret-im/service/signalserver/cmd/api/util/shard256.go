package util

import "crypto/sha256"

func Sha256Str(str string) []byte {

	h := sha256.New()

	h.Write([]byte(str))

	return h.Sum(nil)

}