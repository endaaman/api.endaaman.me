// @APIVersion 1.0.0
// @Title API for endaaman.me
// @Description api.endaaman.me
// @Contact buhibuhidog@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/endaaman/api.endaaman.me/controllers"
	"github.com/astaxie/beego"
)

func init () {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/articles",
			beego.NSInclude(
				&controllers.ArticleController{},
			),
		),
		beego.NSNamespace("/categories",
			beego.NSInclude(
				&controllers.CategoryController{},
			),
		),
		beego.NSNamespace("/sessions",
			beego.NSInclude(
				&controllers.SessionController{},
			),
		),
		beego.NSNamespace("/files",
			beego.NSInclude(
				&controllers.FileController{},
			),
		),
		beego.NSNamespace("/misc",
			beego.NSInclude(
				&controllers.MiscController{},
			),
		),
	)

	beego.ErrorController(&controllers.ErrorController{})
	beego.AddNamespace(ns)
}
