package controllers

import (
	// "mime/multipart"
	// "strings"
	// "encoding/json"
	// "github.com/astaxie/beego/logs"
	"fmt"

	"github.com/endaaman/api.endaaman.me/services"
)

type FileController struct {
	BaseController
}

func (c *FileController) Prepare() {
	c.BaseController.Prepare()
	if !c.IsAdmin {
		c.Respond401()
		c.StopRun()
		return
	}
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
		return
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
	err := services.DeleteFile(rel)
	if err != nil {
		c.Respond400e(err)
		return
	}
}

// @Title Upload file
// @Description delte file
// @Param	files	formData 	true		files
// @Success 200
// @router /* [post]
func (c *FileController) Upload() {
	rel := c.Ctx.Input.Param(":splat")
	isDir := services.IsDir(rel)
	if !isDir {
		c.Respond400(fmt.Sprintf("Target dir `%s` is not directory", rel))
		return
	}

	headers, err := c.GetFiles("files")
	if err != nil {
		c.Respond400f("Uploaded files should be under name `files`: %s", err.Error())
		return
	}

	err = services.SaveFiles(rel, headers)
	if err != nil {
		c.Respond400f("Failed to save files: %s", err.Error())
		return
	}
	c.RespondSimple("success")
}

// @Title Mkdir
// @Description delte file
// @Param	files	formData 	true		files
// @Success 200
// @router /* [put]
func (c *FileController) Mkdir() {
	rel := c.Ctx.Input.Param(":splat")
	err := services.Mkdir(rel)
	if err != nil {
		c.Respond400f("Failed to make dir `%s`: %s", rel, err.Error())
		return
	}
	c.RespondSimple("success")
}

type FileMoveRequest struct {
	Dest string `json:"dest"`
}

// @Title Delete file
// @Description delte file
// @Param	oldName	body 		true	"old name"
// @Param	newName	body 		true	"new name"
// @Success 200
// @router /* [patch]
func (c *FileController) Rename() {
	req := FileMoveRequest{}
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	rel := c.Ctx.Input.Param(":splat")
	err := services.MoveFile(rel, req.Dest)
	if err != nil {
		c.Respond400e(err)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(204)
}
