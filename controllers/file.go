package controllers

import (
	// "fmt"
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
	files, err := services.ListDir(rel)
	if err != nil {
		c.Respond400e(err)
	}
	c.Data["json"] = files
	c.ServeJSON()
}

// @Title Delte file
// @Description delte file
// @Success 200
// @router /* [delete]
func (c *FileController) Delete() {
	rel := c.Ctx.Input.Param(":splat")
	err := services.Delete(rel)
	if err != nil {
		c.Respond400e(err)
		return
	}
}


// @Title Delte file
// @Description delte file
// @Param	oldName	body 		true	"new name"
// @Param	newName	body 		true	"new name"
// @Success 200
// @router /rename [patch]
func (c *FileController) Rename() {
	println("RENAME")
}
