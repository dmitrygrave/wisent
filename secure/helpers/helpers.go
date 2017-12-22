package helpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
)

// GenerateSecureBytes generates a securely random byte slice
func GenerateSecureBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GenerateSecureString generates a securely random base64 encoded string
func GenerateSecureString(n int) (string, error) {
	bytes, err := GenerateSecureBytes(n)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), err
}

// SecureCompare securely compares two byte slices
// NOTE: length should be compared before calling
func SecureCompare(a, b []byte) bool {
	if subtle.ConstantTimeCompare(a, b) == 1 {
		return true
	}

	return false
}

// SecureCompareStrings securely compares two strings
func SecureCompareStrings(a, b string) bool {
	return SecureCompare([]byte(a), []byte(b))
}
