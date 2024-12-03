package randutil

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	idAlphabet           = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	idPocketBaseIDLength = 15
)

// GenerateID generates a random ID with the given length.
//
// It uses idAlphabet as the character set.
func GenerateID(length int) string {
	return gonanoid.MustGenerate(idAlphabet, length)
}

// GenerateIDForPocketBase generates a random ID for PocketBase.
//
// It's length is 15 and it uses idAlphabet as the character set.
//
// With this configuration, it requires 3 Trillion IDs to have a 1%
// probability of at least one collision.
func GenerateIDForPocketBase() string {
	return GenerateID(idPocketBaseIDLength)
}
