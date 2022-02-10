package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

type Crypto interface {
	UserGeneratePassword(password string) string
	UserVerifyPassword(userPwd, userStoredPwd string) bool
}

type crypto struct{}

func NewCrypto() *crypto {
	return &crypto{}
}

func (c crypto) UserGeneratePassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hashedPassword)
}

func (c crypto) UserVerifyPassword(userPwd, userStoredPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userStoredPwd), []byte(userPwd))
	return err == nil
}

//func UserGeneratePassword(password string) string {
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
//	return string(hashedPassword)
//}
//
//func UserVerifyPassword(userPwd, userStoredPwd string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(userStoredPwd), []byte(userPwd))
//	return err == nil
//}
