package common

import (
	"golang.org/x/crypto/bcrypt"
)

// Cost is the cost of the bcrypt digest. This is hard-coded at 10 to match
// the value that Rails/Devise uses. We're porting a Rails database and we
// want to be able to use existing passwords instead of forcing users to
// reset them.
const Cost = 10

// EncryptPassword returns the password run through bcrypt 2a with the
// specified salt and cost.
func EncryptPassword(password string) (string, error) {
	digest, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(digest), nil
}

// ComparePasswords compares a plaintext password against a hashed password,
// returning true if they match, false otherwise.
func ComparePasswords(hashed, plaintext string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext)) == nil
}
