package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

const PassWordCost = 12

func HashPassword(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return ""
	}
	return string(hashPassword)
}

func VerifyHashPassword(password, hashPassword string) bool {
	hashBytes := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(hashBytes, []byte(password))
	return err == nil

}
