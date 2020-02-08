package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func RetrieveArticles(aaCh chan<- []*models.Article) {
	if (infras.IsLoading()) {
		ch := make(chan bool)
		go infras.ReadAllArticles(ch)
		<-ch
	}

    aaCh <- infras.GetCachedArticles()
}
