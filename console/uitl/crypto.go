package util

import "golang.org/x/crypto/bcrypt"

func CryptPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}

func CheckPassword(crtyptPwd string,inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(crtyptPwd), []byte(inputPwd))
	if err != nil {
		return false
	}
	return true
}

