package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/endaaman/api.endaaman.me/controllers"
)

func init() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.AddNamespace(beego.NewNamespace("/v1",
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
	))

	beego.AddNamespace(beego.NewNamespace("/static",
		beego.NSInclude(
			&controllers.StaticController{},
		),
	))
}
