package util

import (
	"math/rand"
	"strconv"
	"time"
)

func GenSalt() (salt string) {
	rand.Seed(time.Now().UnixNano())
	salt = ""
	for i := 1; i < 10; i++ {
		salt = salt + strconv.Itoa(rand.Intn(9))
	}
	return
}
