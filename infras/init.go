package infras

import (
    "github.com/astaxie/beego"
	"github.com/endaaman/api.endaaman.me/utils"
)

func PrepareDirs() {
	ch := make(chan bool)
	go func() {
		utils.EnsureDir(beego.AppConfig.String("private_dir"))
		utils.EnsureDir(beego.AppConfig.String("articles_dir"))
		ch<- true
	}()
	<-ch
}
