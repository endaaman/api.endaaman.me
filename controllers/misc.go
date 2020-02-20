package controllers

import (
	// "encoding/json"
	// "github.com/astaxie/beego"
)

type MiscController struct {
	BaseController
	admin bool
}

func (c *MiscController) Prepare() {
}

// @Title Get warnings
// @Description get warnings
// @Success 200 {string[]} string[]
// @router /warnings [get]
func (c *MiscController) Get() {
	c.Data["json"] = "hi"
	c.ServeJSON()
}
