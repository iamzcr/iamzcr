package admin

import (
	"iamzcr/middleware"
	"crypto/md5"
	"encoding/hex"
)

func generateTestToken() (string, error) {
	return middleware.GenerateToken(1, "test")
}

func validatePassword(password, salt, hash string) bool {
	pwd := password + salt
	h := md5.Sum([]byte(pwd))
	return hex.EncodeToString(h[:]) == hash
}
