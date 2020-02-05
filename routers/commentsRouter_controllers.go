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

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/endaaman/api.endaaman.me/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
