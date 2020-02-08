package controllers

import (
	// "encoding/json"
	"github.com/astaxie/beego"

	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/services"
)

type ArticleController struct {
	beego.Controller
	admin bool
}

func (c *ArticleController) Prepare() {
}

// @Title GetArticles
// @Description get all articles
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (c *ArticleController) Get() {
    aaCh := make(chan []*models.Article)
	go services.RetrieveArticles(aaCh)
	c.Data["json"] = <-aaCh
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
