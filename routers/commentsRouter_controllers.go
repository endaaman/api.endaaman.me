package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ArticleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
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

}
