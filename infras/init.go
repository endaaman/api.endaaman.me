package infras

import (
	"fmt"

	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/utils"
)

func CheckDirs() {
	if !utils.IsDir(config.GetArticlesDir()) {
		panic(fmt.Sprintf("%s is not directory.", config.GetArticlesDir()))
	}

	if !utils.IsDir(config.GetPrivateDir()) {
		panic(fmt.Sprintf("%s is not directory.", config.GetPrivateDir()))
	}
}
