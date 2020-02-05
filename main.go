package main

import (
	// "fmt"
	_ "github.com/endaaman/api.endaaman.me/routers"

    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    orm.RegisterDriver("mysql", orm.DRMySQL)
    orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))

    // err := orm.RunSyncdb("default", false, true)
    // if err != nil {
    //     fmt.Println(err)
    // }
	orm.RunCommand()

    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }
    beego.Run()
}

