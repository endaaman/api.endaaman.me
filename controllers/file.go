package controllers

import (
	"fmt"
	// "mime/multipart"
	// "strings"
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
	if len(headers) == 0 {
		c.Respond400("Uploaded file is empty")
		return
	}

	m := make(map[string]bool)
	for _, header := range headers {
		name := header.Filename
		if m[name] {
			err = fmt.Errorf("Duplicated files(%s) are uploaded", name)
			break
		}
		if !m[name] {
			m[name] = true
		}
		if services.Exists(name) {
			err = fmt.Errorf("The file(%s) already exists.", name)
			break
		}
	}

	if err != nil {
		c.Respond400e(err)
		return
	}

	if len(headers) < 1 {
		c.Respond400("No files uploaded")
		return
	}

	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			c.Respond400f("Failed to open file `%s`:  %v", header.Filename, err)
			return
		}
		err = services.SaveToFile(header.Filename, file)
		if err != nil {
			c.Respond400f("Failed to save file `%s`:  %v", header.Filename, err)
			return
		}
	}
	c.RespondSimple("success")
}

// @Title Delete file
// @Description delte file
// @Param	oldName	body 		true	"old name"
// @Param	newName	body 		true	"new name"
// @Success 200
// @router /rename [patch]
func (c *FileController) Rename() {
	println("RENAME")
}
