package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func RetrieveArticles(ch chan<- []*models.Article) {
	infras.WaitIO()
    ch <- infras.GetCachedArticles()
}

func AppendArticle(a *models.Article, ch chan<- error) {
	infras.WriteArticle(a).Wait()
	infras.AwaitNextChange()
	ch <- nil
}
