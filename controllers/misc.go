package controllers

import (
	// "encoding/json"
	// "github.com/astaxie/beego"
	"github.com/endaaman/api.endaaman.me/services"
)

type MiscController struct {
	BaseController
}

// @Title Get warnings
// @Description get warnings
// @Success 200 {string[]} string[]
// @router /warnings [get]
func (c *MiscController) GetWarnings() {
	c.Data["json"] = services.GetWarnings()
	c.ServeJSON()
}

// @Title Generate hash
// @Param	password	body 		true	"Password"
// @Success 201 Success
// @Failure 400 Validation error
// @router /genhash [post]
func (c *MiscController) GenHash() {
	req := SessionRequest{}
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	hash, err := services.GeneratePasswordHash(req.Password)
	if err != nil {
		c.Respond400e(err)
		return
	}
	c.Data["json"] = map[string]string{"hash": hash}
	c.ServeJSON()
}
