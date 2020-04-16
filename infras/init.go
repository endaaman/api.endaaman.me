package infras

import (
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/utils"
)

func PrepareDirs() {
	err := utils.EnsureDir(config.GetPrivateDir())
	if err != nil {
		panic(err)
	}
	err = utils.EnsureDir(config.GetArticlesDir())
	if err != nil {
		panic(err)
	}
}
