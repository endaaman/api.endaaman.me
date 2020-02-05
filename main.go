package main

import (
	// "regexp"
	_ "github.com/endaaman/api.endaaman.me/routers"
	// "github.com/endaaman/api.endaaman.me/infras"
    "github.com/astaxie/beego"
)

func main() {
	go infras.StartWatching()
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }
    beego.Run()
}
