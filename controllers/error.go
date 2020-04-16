package controllers

import "github.com/astaxie/beego/logs"

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	logs.Debug("No resource found on the URL(%s)", c.Ctx.Input.URI())
	c.RespondSimple("Page not found.")
}
