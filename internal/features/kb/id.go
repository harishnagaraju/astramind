package kb

import (
	"crypto/rand"
	"encoding/hex"
)

func generateDocumentID() string {
	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
