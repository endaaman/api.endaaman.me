package services

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

// never use this cuz the admin must be just me
func GeneratePasswordHash(password string)(string,error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func ValidatePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateToken(expiration uint) (string, error) {
	secret := beego.AppConfig.String("secret_key_base")
	now := time.Now()
	exp := now.Add(time.Duration(expiration) * 24 * time.Hour)

	jsonToken := paseto.JSONToken{
		// Audience:   "",
		Issuer:     "api.endaaman.me",
		// Jti:        "",
		// Subject:    "",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}
	jsonToken.Set("", "")

	v2 := paseto.NewV2()
	token, err := v2.Encrypt([]byte(secret), jsonToken, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(token string) bool {
	secret := beego.AppConfig.String("secret_key_base")
	var newJsonToken paseto.JSONToken
	var newFooter string

	v2 := paseto.NewV2()
	err := v2.Decrypt(token, []byte(secret), &newJsonToken, &newFooter)
	if err != nil {
		return false
	}
	return true
}
