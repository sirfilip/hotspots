package services

import "golang.org/x/crypto/bcrypt"

type Crypter interface {
	Hash(string) (string, error)
	Check(hashed, secret string) error
}

type Bcrypt struct {
	cost int
}

func NewBcrypt() Bcrypt {
	return Bcrypt{cost: bcrypt.DefaultCost}
}

func (crypt Bcrypt) Hash(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), crypt.cost)
	return string(hash), err
}

func (crypt Bcrypt) Check(hashed, secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(secret))
}
