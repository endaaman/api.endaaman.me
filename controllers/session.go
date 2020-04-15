package controllers

import (
	"github.com/endaaman/api.endaaman.me/services"
)

type SessionController struct {
	BaseController
	admin bool
}

type SessionResponse struct {
	Token string `json:"token"`
}

type SessionRequest struct {
	Password string `json:"password"`
}

// @Title Check if authenticated
// @Success 200 You are me
// @Success 401 You are not me
// @router / [get]
func (c *SessionController) Check() {
	var m string
	if c.IsAdmin {
		m = "You are me."
	} else {
		m = "You are not me."
		c.Ctx.Output.SetStatus(401)
	}
	c.Data["json"] = NewSimpleResponse(m)
	c.ServeJSON()
}

// @Title Create session
// @Param	password	body	true	"Password"
// @Success 201 Success
// @Failure 400 Validation error
// @router / [post]
func (c *SessionController) Login() {
	req := SessionRequest{}
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	suc := services.ValidatePassword(req.Password)
	if !suc {
		c.Respond401()
		return
	}

	token, err := services.GenerateToken(7)
	if err != nil {
		c.Respond400(err.Error())
		return
	}
	res := &SessionResponse{Token: token}
	c.Data["json"] = &res
	c.ServeJSON()
}

// @Title Renew token
// @Success 201 Success
// @Failure 400 Validation error
// @router /renew [post]
func (c *SessionController) Renew() {
	c.Respond400("not implemented")
}
