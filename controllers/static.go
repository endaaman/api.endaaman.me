package controllers

import (
	"io/ioutil"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/endaaman/api.endaaman.me/config"
)

type StaticController struct {
	BaseController
}

// @Title Serve static file
// @Success 200 You are me
// @Success 401 If you are not me
// @router /* [get]
func (c *StaticController) Get() {
	restrictedDirs := []string{
		config.GetArticlesDirname(),
		config.GetPrivateDirname(),
	}

	rel := c.Ctx.Input.Param(":splat")

	splitted := strings.SplitN(rel, "/", 2)
	base := splitted[0]
	for _, restricted := range restrictedDirs {
		if restricted == base {
			if !c.IsAdmin {
				c.Respond401()
				return
			}
		}
	}

	if config.IsDev() {
		baseDir := config.GetSharedDir()
		buf, err := ioutil.ReadFile(filepath.Join(baseDir, rel))
		if err != nil {
			c.Respond400e(err)
			return
		}
		c.Ctx.Output.ContentType(rel)
		c.Ctx.Output.Body(buf)
		return
	}

	u, err := url.Parse("http://localhost:3002/" + rel)
	if err != nil {
		c.Respond400e(err)
		return
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	req := c.Ctx.Request.Clone(c.Ctx.Request.Context())
	// modify path: HOST/static/hoge/fuga -> HOST/hoge/fuga
	req.URL.Path = "/" + rel
	rp.ServeHTTP(c.Ctx.ResponseWriter, req)
}
