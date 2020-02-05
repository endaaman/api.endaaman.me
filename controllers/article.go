package controllers

import (
	"github.com/endaaman/api.endaaman.me/models"
	// "encoding/json"
	"github.com/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) Prepare() {
	print("pre")
}

// @Title GetAllArticles
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (c *ArticleController) Get() {
	c.Data["json"] = models.GetAllArticles()
	c.ServeJSON()

}

// @Title CreateArticle
// @Description create object
// @Param	article	body 	models.Object	true		"The article content"
// @Success 201 {body} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (c *ArticleController) Post() {
	// var a models.Article
	// json.Unmarshal(o.Ctx.Input.RequestBody, &a)
	// objectid := models.AddOne(ob)
	// o.Data["json"] = map[string]string{"ObjectId": objectid}
	c.ServeJSON()
}
