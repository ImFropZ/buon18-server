package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPwd(rawPwd string) (string, error) {
	rawPwdBytes := []byte(rawPwd)
	hashedPwd, err := bcrypt.GenerateFromPassword(rawPwdBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func ComparePwd(rawPwd, hashedPwd string) bool {
	rawPwdBytes := []byte(rawPwd)
	hashedPwdBytes := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(hashedPwdBytes, rawPwdBytes)
	return err == nil
}
