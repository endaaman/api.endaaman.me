package services

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/utils"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

const PASSWORD_HASH_FILE = "password_hash"
const SECRET_KEY_BASE_FILE = "secret_key_base"

func init() {
	_ = getSecretKey()
	_ = getPasswordHash()
}

func getPasswordHash() string {
	hashPath := filepath.Join(config.GetPrivateDir(), PASSWORD_HASH_FILE)
	if utils.FileExists(hashPath) {
		buf, err := ioutil.ReadFile(hashPath)
		if err != nil {
			return ""
		}
		splitted := strings.SplitN(string(buf), "\n", 2)
		return splitted[0]
	}

	logs.Warn("Password hash file(%s) does not exist. Loaded from conf instead", hashPath)
	hash, err := GeneratePasswordHash(config.GetPassword())
	if err != nil {
		panic(err)
	}
	return hash
}

func getSecretKey() string {
	secretPath := filepath.Join(config.GetPrivateDir(), SECRET_KEY_BASE_FILE)
	if utils.FileExists(secretPath) {
		buf, err := ioutil.ReadFile(secretPath)
		if err != nil {
			panic(err)
		}
		splitted := strings.SplitN(string(buf), "\n", 2)
		return splitted[0]
	}

	logs.Warn("Secret key base file(%s) does not exist. Loaded from conf instead", secretPath)
	return config.GetSecretKeyBase()
}

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func ValidatePassword(password string) bool {
	hash := getPasswordHash()
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateToken(expiration uint) (string, error) {
	secret := getSecretKey()
	now := time.Now()
	// 1 month
	exp := now.Add(time.Duration(expiration) * 30 * 24 * time.Hour)

	jsonToken := paseto.JSONToken{
		// Audience:   "",
		Issuer: "api.endaaman.me",
		// Jti:        "",
		// Subject:    "",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}

	v2 := paseto.NewV2()
	token, err := v2.Encrypt([]byte(secret), jsonToken, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(token string) bool {
	secret := getSecretKey()
	var newJsonToken paseto.JSONToken
	var newFooter string

	v2 := paseto.NewV2()
	err := v2.Decrypt(token, []byte(secret), &newJsonToken, &newFooter)
	if err != nil {
		return false
	}
	return true
}
