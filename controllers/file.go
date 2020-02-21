package controllers

import (
	"fmt"
	// "encoding/json"
	// "github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/services"
)

type FileController struct {
	BaseController
}

// @Title Get files
// @Description get files
// @Success 200 []
// @router /* [get]
func (c *FileController) ListDir() {
	rel := c.Ctx.Input.Param(":splat")
	if !services.IsDir(rel) {
		c.Respond400(fmt.Sprintf("Can not read the path: `%s`", rel))
		return
	}
	c.Data["json"] = services.ListDir(rel)
	c.ServeJSON()
}

// @Title Delte file
// @Description delte file
// @Success 200
// @router /* [delete]
func (c *FileController) Delete() {
	rel := c.Ctx.Input.Param(":splat")
	if !services.IsDir(rel) {
		c.Respond400(fmt.Sprintf("Can not read the path: `%s`", rel))
		return
	}
	c.Data["json"] = "DEL"
	c.ServeJSON()
}
