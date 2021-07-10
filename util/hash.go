package util

import "golang.org/x/crypto/bcrypt"

// Hash return hash result from string input
func Hash(plainString string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainString), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// HashCompare compare hashed string with plain string
func HashCompare(plainString string, hashedString string) bool {
	hashedBytes := []byte(hashedString)
	plainBytes := []byte(plainString)

	if err := bcrypt.CompareHashAndPassword(hashedBytes, plainBytes); err == nil {
		return true
	}

	return false
}
