package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const saltSize = 16
const bcryptCost = 16

func GenerateRandomSalt() (string, error) {
	salt := make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	return hex.EncodeToString(salt), err
}

func HashPasswordWithSalt(password string, salt string) (string, error) {

	passwordBytes := []byte(password + salt) // value and salt

	hashBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcryptCost)

	return string(hashBytes), err
}
