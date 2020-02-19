package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Create",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/:category/:slug`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Remove",
            Router: `/:category/:slug`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:MiscController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:MiscController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/warnings`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"],
        beego.ControllerComments{
            Method: "Check",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:SessionController"],
        beego.ControllerComments{
            Method: "Renew",
            Router: `/renew`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
