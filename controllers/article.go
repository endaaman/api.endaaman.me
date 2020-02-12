package controllers

import (
	// "net/url"
	"encoding/json"
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

func (c *ArticleController) ServeJSONText(data []byte) {
	c.Ctx.Output.Header("Content-Description", "File Transfer")
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	// c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")
	c.Ctx.Output.Header("Expires", "0")
	c.Ctx.Output.Header("Cache-Control", "must-revalidate")
	c.Ctx.Output.Header("Pragma", "public")
	c.Ctx.Output.Body(data)
}

// @Title Get all articles
// @Success 200 {object} models.Article
// @router / [get]
func (c *ArticleController) Get() {
    ch := make(chan []*models.Article)
	go services.RetrieveArticles(ch)
	c.Data["json"] = <-ch
	c.ServeJSON()
}

// @Title Create an article
// @Param	article	body 	models.Article	true	"The article content"
// @Success 201 Success
// @Failure 400 Validation error
// @router / [post]
func (c *ArticleController) Post() {
	a := models.NewArticle()
	err := json.Unmarshal(c.Ctx.Input.RequestBody, a)
	if err != nil {
		c.Data["json"] = map[string]string{"message": err.Error()}
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	messages := a.Validate()
	if messages != nil {
		c.Data["json"] = map[string]interface{}{"errors": messages, "message": "There are some value errors."}
		// s, _ := json.MarshalIndent(a, "", "  ")
		// fmt.Println(string(s))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	ch := make(chan error)
	go services.AppendArticle(a, ch)
	err = <-ch
	if err != nil {
		c.Data["json"] = err
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}
    ch2 := make(chan []*models.Article)
	go services.RetrieveArticles(ch2)
	c.Data["json"] = <-ch2
	c.ServeJSON()
}
