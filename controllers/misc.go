package controllers

import (
	// "encoding/json"
	// "github.com/astaxie/beego"
	"github.com/endaaman/api.endaaman.me/services"
)

type MiscController struct {
	BaseController
}

func (c *MiscController) Prepare() {
	c.BaseController.Prepare()
	if !c.IsAdmin {
		c.Respond401()
		c.StopRun()
		return
	}
}

type statusAnnotation struct {
	Warnings map[string][]string `json:"warnings"`
	Watcher  struct {
		IsActive  bool   `json:"isActive"`
		LastError string `json:"lastError"`
	} `json:"watcher"`
}

// @Title Get status
// @Description get status
// @Success 200
// @router /status [get]
func (c *MiscController) GetStatus() {
	status := statusAnnotation{}
	status.Warnings = services.GetWarnings()
	status.Watcher.IsActive = services.IsWatcherActive()
	status.Watcher.LastError = services.GetWathcerLastError()
	c.Data["json"] = status
	c.ServeJSON()
}

// @Title Generate hash
// @Param	password	body 		true	"Password"
// @Success 201 Success
// @Failure 400 Validation error
// @router /genhash [post]
func (c *MiscController) GenHash() {
	req := SessionRequest{}
	if !c.ExpectJSON(&req) {
		c.Respond400InvalidJSON()
		return
	}

	hash, err := services.GeneratePasswordHash(req.Password)
	if err != nil {
		c.Respond400e(err)
		return
	}
	c.Data["json"] = map[string]string{"hash": hash}
	c.ServeJSON()
}

// @Title Restart watcher
// @Success 201 Success
// @Failure 400 Watcher already started
// @router /watcher/restart [post]
func (c *MiscController) RestartWatching() {
	err := services.RestartWatcher()
	if err != nil {
		c.Respond400e(err)
		return
	}
	c.RespondSimple("Queued watcher restart")
}
