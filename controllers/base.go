package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/services"
)

type BaseController struct {
	beego.Controller
	IsAdmin bool
}

type SimpleResponse struct {
	Message string `json:"message"`
}

type ValidationFailureResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

const TOKEN_PREFIX = "Bearer"
const BYPASS_PARAM = "x"

func NewSimpleResponse(message string) *SimpleResponse {
	p := SimpleResponse{}
	p.Message = message
	return &p
}

func (c *BaseController) RespondSimple(message string) {
	c.Data["json"] = NewSimpleResponse(message)
	c.ServeJSON()
}

func (c *BaseController) Respond404() {
	c.Data["json"] = NewSimpleResponse("The resource does not exist.")
	c.Ctx.ResponseWriter.WriteHeader(400)
	c.ServeJSON()
}

func (c *BaseController) Respond400(message string) {
	logs.Warn("[400] - %s", message)
	c.Data["json"] = NewSimpleResponse(message)
	c.Ctx.ResponseWriter.WriteHeader(400)
	c.ServeJSON()
}

func (c *BaseController) Respond400f(format string, args ...interface{}) {
	c.Respond400(fmt.Sprintf(format, args...))
}

func (c *BaseController) Respond400e(err error) {
	c.Respond400(err.Error())
}

func (c *BaseController) Respond401() {
	c.Data["json"] = NewSimpleResponse("You are not me.")
	c.Ctx.ResponseWriter.WriteHeader(401)
	c.ServeJSON()
}

func (c *BaseController) Respond403() {
	c.Data["json"] = NewSimpleResponse("You are not true me.")
	c.Ctx.ResponseWriter.WriteHeader(403)
	c.ServeJSON()
}

func (c *BaseController) Respond400InvalidJSON() {
	c.Respond400("Invalid JSON format.")
}

func (c *BaseController) Respond400ValidationFailure(err *models.ValidationError) {
	res := ValidationFailureResponse{
		Message: "There are some value errors.",
		Errors:  err.Messages,
	}
	c.Data["json"] = res
	c.Ctx.ResponseWriter.WriteHeader(400)
	c.ServeJSON()
}

func (c *BaseController) ExpectJSON(data interface{}) bool {
	err := json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		logs.Warn("JSON parse error:", err.Error())
		logs.Warn(string(c.Ctx.Input.RequestBody))
		return false
	}
	return true
}

func (c *BaseController) authorizeToken(token string) {
	c.IsAdmin = services.ValidateToken(token)
	if !c.IsAdmin {
		logs.Warn("Tried to authenticate invalid token.")
	}
}

func (c *BaseController) authorize() {
	// dev and ?x
	if beego.BConfig.RunMode == "dev" {
		_, bypass := c.Ctx.Request.URL.Query()[BYPASS_PARAM]
		if bypass {
			c.IsAdmin = true
			logs.Warn("Bypassed to admin for development")
			return
		}
	}

	// auth by 'Authorization' header
	var token string
	rawToken := c.Ctx.Input.Header("Authorization")
	splitted := strings.SplitN(rawToken, " ", 2)
	if len(splitted) == 2 && splitted[0] == TOKEN_PREFIX {
		c.authorizeToken(splitted[1])
		return
	}

	// auth by 'token' cookie
	token = c.Ctx.GetCookie("token")
	if token != "" {
		c.authorizeToken(token)
	}
}

func (c *BaseController) Prepare() {
	c.authorize()
}
