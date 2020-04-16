package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/infras"
	_ "github.com/endaaman/api.endaaman.me/routers"
	"github.com/endaaman/api.endaaman.me/services"
)

//go:generate sh -c "echo 'package routers; import \"github.com/astaxie/beego\"; func init() {beego.BConfig.RunMode = beego.DEV}' > routers/0.go"
//go:generate sh -c "echo 'package routers; import \"os\"; func init() {os.Exit(0)}' > routers/z.go"
//go:generate go run $GOFILE
//go:generate sh -c "rm routers/0.go routers/z.go"

func main() {
	infras.PrepareDirs()
	services.ReadAllArticles()
	go infras.StartWatching()
	if config.IsDev() {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}), true)

	beego.Run()
}
