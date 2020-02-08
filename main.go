package main

import (
	// "fmt"
	// "log"
	// "time"
	// "regexp"
	// "github.com/radovskyb/watcher"
	// "github.com/bep/debounce"
    "github.com/astaxie/beego"
	_ "github.com/endaaman/api.endaaman.me/routers"
	"github.com/endaaman/api.endaaman.me/infras"
)

func main() {
	go infras.StartWatching()
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

    }
    beego.Run()
}
