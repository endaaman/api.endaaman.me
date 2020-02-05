// +build wireinject
package main

import (
    "github.com/google/wire"
	"github.com/endaaman/api.endaaman.me/routers"
	"github.com/endaaman/api.endaaman.me/infras"
	"github.com/endaaman/api.endaaman.me/usecases"
	"github.com/endaaman/api.endaaman.me/controllers"
)

func Inject() *routers.Router {
    wire.Build(
		infras.NewArticleRepository,
		usecases.NewArticleUsecase,
		controllers.NewArticleController,
		routers.RegisterRouter,
    )
    return nil
}
