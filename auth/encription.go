package auth

import "golang.org/x/crypto/bcrypt"

func Encode(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func EncodeStr(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Compare(password []byte, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
