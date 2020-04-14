package controllers

import "fmt"

type StaticController struct {
	BaseController
}

// @Title Serve static file
// @Success 200 You are me
// @Success 401 You are not me
// @router / [get]
func (c *StaticController) Get() {
	fmt.Println("static")
	c.Data["json"] = "res"
	c.ServeJSON()
}
