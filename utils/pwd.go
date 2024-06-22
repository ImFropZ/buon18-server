package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPwd(rawPwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(rawPwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func ComparePwd(rawPwd, hashedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(rawPwd))
	return err == nil
}
