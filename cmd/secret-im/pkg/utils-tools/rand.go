package utils

import (
	"crypto/rand"
	"math/big"
)

// 安全随机int值
func SecureRandInt() int {
	return SecureRandIntN(IntMax)
}

func SecureRandIntN(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

// 安全随机int64值
func SecureRandInt64() int64 {
	return SecureRandInt64N(Int64Max)
}

func SecureRandInt64N(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return n.Int64()
}
