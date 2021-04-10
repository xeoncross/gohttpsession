package sessiontoken

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encode bytes to a URL safe, RFC 4648, base64 string
func Encode(s []byte) string {
	return base64.URLEncoding.EncodeToString(s)
}

// Decode a URL safe, RFC 4648, base64 string to raw bytes
// On failure, returns nil
func Decode(s string) []byte {
	b, err := base64.URLEncoding.DecodeString(s)
	if err == nil && len(b) != 0 {
		return b
	}
	return nil
}

// New creates a random token with the given length in bytes.
// On failure, returns nil.
//
// Callers should explicitly check for the possibility of a nil return, treat
// it as a failure of the system random number generator, and not continue.
func New(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
