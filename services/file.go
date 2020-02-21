package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func ListDir(rel string) []*models.File {
	ch := make(chan []*models.File)
    go infras.ListDir(rel, ch)
	return <-ch
}

func IsDir(rel string) bool {
	ch := make(chan bool)
    go infras.IsDir(rel, ch)
	return <-ch
}
