package ops

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(parts ...[]byte) string {
	hasher := sha256.New()
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		_, _ = hasher.Write(part)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
