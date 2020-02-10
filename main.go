package main

import (
	"os"
    "github.com/astaxie/beego"
	_ "github.com/endaaman/api.endaaman.me/routers"
	"github.com/endaaman/api.endaaman.me/infras"
    "github.com/astaxie/beego/logs"
)


func main() {
    secret := os.Getenv("SECRET_KEY_BASE")
	if secret == "" {
		secret = "THIS IS DEGEROUS SECRET KEY BASE"
		logs.Warn("using degerous key")
	}
	beego.AppConfig.Set("secret", secret)

	infras.ReadAllArticles().Wait()
	go infras.StartWatching()
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

    }
    beego.Run()
}
