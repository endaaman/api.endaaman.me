package controllers

import (
	"github.com/endaaman/api.endaaman.me/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	// var user models.User
	// json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	// uid := models.AddUser(user)
	c.Data["json"] = map[string]string{"uid": "hoge"}
	c.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (c *UserController) GetAll() {
	users := models.GetAllUsers()
	c.Data["json"] = users
	c.ServeJSON()
}

// @Title GetOne
// @Description get user by id
// @Param	id		path 	string	true		"The uid you want to update"
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) Get() {
	id, err := c.GetInt(":id")

	if err != nil {
		c.Data["json"] = err.Error()
	}

	user, err := models.GetUser(id)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	if err == nil {
		c.Data["json"] = user
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	id		path 	string	true		"The uid you want to update"
// @Param	body	body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:id [put]
func (c *UserController) Put() {
	id := c.GetString(":id")
	if id != "" {
		var user models.User
		json.Unmarshal(c.Ctx.Input.RequestBody, &user)
		u2, err := models.UpdateUser(id, &user)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = u2
		}
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	id	path 	string	true		"id for user to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (c *UserController) Delete() {
	id := c.GetString(":id")
	models.DeleteUser(id)
	c.Data["json"] = "delete success!"
	c.ServeJSON()
}
