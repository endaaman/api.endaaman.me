package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func GetCategorys() []*models.Category {
	infras.WaitIO()
    return infras.GetCachedCategorys()
}
