package services

import (
	"fmt"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func RetrieveArticles(ch chan<- []*models.Article) {
	infras.WaitReader()
    ch <- infras.GetCachedArticles()
}

func AppendArticle(a *models.Article, ch chan<- error) {
	// save article
	fmt.Println(a.ToText())
	infras.WaitReader()
    ch <- nil
}
