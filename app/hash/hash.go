package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Secret struct {
	hash string
}

const hashCost = 10

var errHash = errors.New("error hashing")

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
