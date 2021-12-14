package hash

import "golang.org/x/crypto/bcrypt"

type Secret struct {
	hash string
}

const hashCost = 10

func NewHash(secret string) (Secret, error) {
	s := Secret{}

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), hashCost)

	if err != nil {

		return s, err
	}

	s.hash = string(hash)

	return s, nil
}

func (s Secret) Compare(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.hash), []byte(password))

	return err == nil
}

func (s Secret) Value() string {
	return s.hash
}
