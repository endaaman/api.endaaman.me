package controllers

import (
	"fmt"
	"encoding/json"
	"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/services"
	"github.com/endaaman/api.endaaman.me/models"
)

type BaseController struct {
	beego.Controller
	IsAdmin bool
}

type SimpleResponse struct {
	Message string `json:"message"`
}

type ValidationFailureResponse struct {
	Message string `json:"message"`
	Errors map[string][]string `json:"errors"`
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
		Errors: err.Messages,
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

func (c *BaseController) Prepare() {
	rawToken := c.Ctx.Input.Header("Authorization")
	splitted := strings.SplitN(rawToken, " ", 2)
	var token string
	if len(splitted) == 2 && splitted[0] == TOKEN_PREFIX {
		token = splitted[1]
		c.IsAdmin = services.ValidateToken(token)
		if c.IsAdmin {
			logs.Info("Successfuly logged in.")
		} else {
			logs.Warn("Tried to authenticate invalid token.")
		}
	} else {
		if beego.BConfig.RunMode == "dev" {
			_, bypass := c.Ctx.Request.URL.Query()[BYPASS_PARAM]
			if bypass {
				c.IsAdmin = true
				logs.Warn("Bypassed to admin for development")
			}
		}
	}
}
