package config

import (
	"fmt"
	"path/filepath"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
)

func getStringValue(key string) string {
	s, err := config.String(key)
	if err != nil {
		panic(err)
	}
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
