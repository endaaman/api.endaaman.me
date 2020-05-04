package controllers

import (
	// "fmt"
	// "net/url"
	// "github.com/astaxie/beego"

	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/services"
)

type ArticleController struct {
	BaseController
}

// @Title Get all articles
// @Success 200 {object} models.Article
// @router / [get]
func (c *ArticleController) GetAll() {
	c.Data["json"] = services.GetArticles(c.IsAdmin)
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

	a := models.NewArticle()
	if !c.ExpectJSON(&a) {
		c.Respond400InvalidJSON()
		return
	}

	err := a.Validate()
	if err != nil {
		if e, ok := err.(*models.ValidationError); ok {
			c.Respond400ValidationFailure(e)
		} else {
			c.Respond400e(err)
		}
		return
	}

	err = services.AddArticle(a)
	if err != nil {
		c.Respond400e(err)
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

	newA := models.NewArticle()
	if !c.ExpectJSON(&newA) {
		c.Respond400InvalidJSON()
		return
	}

	err := newA.Validate()
	if err != nil {
		if e, ok := err.(*models.ValidationError); ok {
			c.Respond400ValidationFailure(e)
		} else {
			c.Respond400e(err)
		}
		return
	}

	err = services.UpdateArticle(oldA, newA)
	if err != nil {
		c.Respond400e(err)
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
		c.Respond400e(err)
		return
	}
	c.Ctx.Output.SetStatus(200)
}
