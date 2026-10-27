package util

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateMD5Hash(input string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(input))
	hashBytes := md5Hash.Sum(nil)
	hash := hex.EncodeToString(hashBytes)
	return hash
}

// HashPassword function : password argument passed is already in MD5.
func HashPassword(password string) (string, error) {
	// salt := config.Get().Encrypt.KeySalt

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// HashPasswordWithMD5 function : password argument passed is not in MD5 yet,
// it is pure string, then need to convert to MD5 first.
func HashPasswordWithMD5(password string) (string, error) {
	// salt := config.Get().Encrypt.KeySalt
	password = GenerateMD5Hash(password)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ComparePassword function : passwordInputted argument passed is already in MD5.
func ComparePassword(passwordInputted, passwordCompared string) bool {
	// salt := config.Get().Encrypt.KeySalt
	err := bcrypt.CompareHashAndPassword([]byte(passwordCompared), []byte(passwordInputted))
	return err == nil
}

// ComparePasswordWithMD5 function : passwordInputted argument passed is not in MD5 yet,
// it is pure string, then need to convert to MD5 first.
func ComparePasswordWithMD5(passwordInputted, passwordCompared string) bool {
	// salt := config.Get().Encrypt.KeySalt

	passwordInputted = GenerateMD5Hash(passwordInputted)
	err := bcrypt.CompareHashAndPassword([]byte(passwordCompared), []byte(passwordInputted))
	return err == nil
}
