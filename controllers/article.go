package controllers

import (
	// "fmt"
	// "net/url"
	// "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/services"
)


type ArticleRequest struct {
	models.Article
}

type ArticleController struct {
	BaseController
	admin bool
}

func NewArticleRequest() *ArticleRequest {
	r := ArticleRequest{}
	r.Article = *models.NewArticle()
	return &r
}

// @Title Get all articles
// @Success 200 {object} models.Article
// @router / [get]
func (c *ArticleController) GetAll() {
	c.Data["json"] = services.GetArticles()
	c.ServeJSON()
}

// @Title Create an article
// @Param	article	body 	models.Article	true	"The article content"
// @Success 201 Success
// @Failure 400 Validation error
// @Failure 401 Auth error
// @router / [post]
func (c *ArticleController) Create() {
	if !c.IsAdmin {
		c.Respond401()
		return
	}

	req := NewArticleRequest()
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	a := &req.Article
	messages := a.Validate()
	if messages != nil {
		c.Respond400ValidationFailure(messages)
		return
	}

	err := services.AddArticle(a)
	if err != nil {
		c.Respond400(err.Error())
		return
	}
	c.Data["json"] = services.IdentifyArticle(a)
	c.ServeJSON()
}

// @Title Update the article
// @Param	article	body 	models.Article	true	"The article content"
// @Success 200 Success
// @Failure 400 Validation error
// @Failure 401 Auth error
// @router /:category/:slug [patch]
func (c *ArticleController) Update() {
	if !c.IsAdmin {
		c.Respond401()
		return
	}

	req := &ArticleRequest{}
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}


	// _, bypass := c.Ctx.Request.URL.Query()[BYPASS_PARAM]
	// if bypass {
	// 	c.IsAdmin = true
	// 	logs.Warn("Bypassed to admin for development")
	// }
	needleCategory := c.Ctx.Input.Param(":category")
	if needleCategory == "-" {
		needleCategory = ""
	}
	needleSlug := c.Ctx.Input.Param(":slug")
	logs.Info("Category: `%s` Slug: `%s`", needleCategory, needleSlug)

	a := services.FindArticle(needleCategory, needleSlug)

	c.Data["json"] = a
	c.ServeJSON()

	// a := &req.Article
	// messages := a.Validate()
	// if messages != nil {
	// 	c.Respond400ValidationFailure(messages)
	// 	return
	// }
}
