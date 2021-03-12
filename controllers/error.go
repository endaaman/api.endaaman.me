package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	logs.Debug("No resource found on the URL(%s)", c.Ctx.Input.URI())
	c.RespondSimple("Page not found.")
}
