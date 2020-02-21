package controllers

import (
	// "fmt"
	// "net/url"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/logs"

	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/services"
)


type ArticleRequest struct {
	models.Article
}

type CategoryController struct {
	BaseController
}

// @Title Get all categorys
// @Success 200 {object} models.Article
// @router / [get]
func (c *CategoryController) GetAll() {
	c.Data["json"] = services.GetCategorys()
	c.ServeJSON()
}
