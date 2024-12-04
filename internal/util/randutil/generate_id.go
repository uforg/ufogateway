package randutil

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	// idAlphabet is the alphabet supported by PocketBase for ID generation.
	idAlphabet           = "0123456789abcdefghijklmnopqrstuvwxyz"
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
// With this configuration, it requires 66 Billion IDs to have a 1%
// probability of at least one collision.
func GenerateIDForPocketBase() string {
	return GenerateID(idPocketBaseIDLength)
}
