package store

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(byteLen int) string {
	b := make([]byte, byteLen)
	rand.Read(b)
	return hex.EncodeToString(b)
}
