package controllers

import (
	"fmt"
	// "net/url"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/logs"

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
// @router /:category/:slug [put]
func (c *ArticleController) Update() {
	if !c.IsAdmin {
		c.Respond401()
		return
	}

	oldA := services.FindArticle(c.Ctx.Input.Param(":category"), c.Ctx.Input.Param(":slug"))
	if oldA == nil {
		c.Respond404()
		return
	}

	req := NewArticleRequest()
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	newA := &req.Article
	fmt.Printf("%+v", newA)
	messages := newA.Validate()
	if messages != nil {
		c.Respond400ValidationFailure(messages)
		return
	}

	err := services.ReplaceArticle(oldA, newA)
	if err != nil {
		c.Respond400(err.Error())
		return
	}
	c.Data["json"] = services.IdentifyArticle(newA)
	c.ServeJSON()
}


// @Title Remove the article
// @Success 200 Success
// @Failure 400 Validation error
// @Failure 401 Auth error
// @router /:category/:slug [delete]
func (c *ArticleController) Remove() {
	if !c.IsAdmin {
		c.Respond401()
		return
	}

	oldA := services.FindArticle(c.Ctx.Input.Param(":category"), c.Ctx.Input.Param(":slug"))
	if oldA == nil {
		c.Respond404()
		return
	}

	err := services.RemoveArticle(oldA)
	if err != nil {
		c.Respond400(err.Error())
		return
	}
	c.Ctx.Output.SetStatus(200)
}
