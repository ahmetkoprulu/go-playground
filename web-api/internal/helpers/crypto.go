package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string, key string) string {
	hash := sha256.New()
	hash.Write([]byte(password + key))
	return hex.EncodeToString(hash.Sum(nil))
}
