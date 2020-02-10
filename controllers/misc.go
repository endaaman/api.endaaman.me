package controllers

import (
	// "encoding/json"
	"github.com/astaxie/beego"
	"github.com/endaaman/api.endaaman.me/services"
)

type MiscController struct {
	beego.Controller
	admin bool
}

func (c *MiscController) Prepare() {
}

// @Title Get warnings
// @Description get warnings
// @Success 200 {string[]} string[]
// @router /warnings [get]
func (c *MiscController) Get() {
    ch := make(chan []string)
	go services.RetrieveWarnings(ch)
	c.Data["json"] = <- ch
	c.ServeJSON()
}
