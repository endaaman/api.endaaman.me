package controllers

import (
    // "github.com/astaxie/beego"
)

type ErrorController struct {
    BaseController
}

func (c *ErrorController) Error404() {
	c.RespondSimple("Page not found.")
}
