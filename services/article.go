package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func searchArticle(aa []*models.Article, category, slug string) *models.Article {
	for _, a := range aa {
		if a.Category == category && a.Slug == slug {
			return a
		}
	}
	return nil
}

func GetArticles() []*models.Article {
	infras.WaitIO()
    return infras.GetCachedArticles()
}

func FindArticle(category, slug string) *models.Article {
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	return searchArticle(aa, category, slug)
}

func AddArticle(a *models.Article) error {
	infras.WaitIO()
	ch := make(chan error)
	infras.WriteArticle(a, ch)
	err := <-ch
	if err != nil {
		return err
	}
	infras.AwaitNextChange()
	return nil
}

func IdentifyArticle(a *models.Article) *models.Article {
	aa := infras.GetCachedArticles()
	return searchArticle(aa, a.Category, a.Slug)
}
