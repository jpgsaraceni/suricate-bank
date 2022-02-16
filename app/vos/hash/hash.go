package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Secret struct {
	hash string
}

const hashCost = 10

var (
	errHash      = errors.New("error hashing")
	errScan      = errors.New("scan failed")
	errScanEmpty = errors.New("scan returned empty")
)

func NewHash(inputSecret string) (Secret, error) {
	s := Secret{}

	hash, err := bcrypt.GenerateFromPassword([]byte(inputSecret), hashCost)
	if err != nil {
		return s, errHash
	}

	s.hash = string(hash)

	return s, nil
}

func (s Secret) Compare(inputHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.hash), []byte(inputHash))

	return err == nil
}

func (s Secret) Value() string {
	return s.hash
}

// Scan implements database/sql/driver Scanner interface.
// Scan parses a string value to hash (if valid) or returns error.
func (s *Secret) Scan(value interface{}) error {
	if value == nil {
		*s = Secret{}

		return errScanEmpty
	}
	if value, ok := value.(string); ok {
		if _, err := bcrypt.Cost([]byte(value)); err != nil { // check if secret is a valid hash
			return errScan
		}

		*s = Parse(value)

		return nil
	}

	return errScan
}

func Parse(secret string) Secret {
	return Secret{hash: secret}
}
