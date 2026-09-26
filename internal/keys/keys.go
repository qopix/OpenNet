package keys

import (
	"crypto/rand"
	"encoding/hex"
)

func Generate() (string, error) {
	data := make([]byte, 32)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}