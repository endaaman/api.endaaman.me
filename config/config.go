package config

import (
	"fmt"
	"path/filepath"

	"github.com/astaxie/beego"
)

func getStringValue(key string) string {
	s := beego.AppConfig.String(key)
	if s == "" {
		panic(fmt.Errorf("Value of key `%s` is empty", key))
	}
	return s
}

func GetSecretKeyBase() string {
	return getStringValue("secret_key_base")
}

func GetPassword() string {
	return getStringValue("password")
}

func GetSharedDir() string {
	return getStringValue("shared_dir")
}

func GetArticlesDirname() string {
	return getStringValue("articles_dirname")
}

func GetPrivateDirname() string {
	return getStringValue("private_dirname")
}

func GetArticlesDir() string {
	sharedPath := getStringValue("shared_dir")
	articlesDir := getStringValue("articles_dirname")
	return filepath.Join(sharedPath, articlesDir)
}

func GetPrivateDir() string {
	sharedPath := getStringValue("shared_dir")
	privateDir := getStringValue("private_dirname")
	return filepath.Join(sharedPath, privateDir)
}

func IsDev() bool {
	return beego.BConfig.RunMode == beego.DEV
}
