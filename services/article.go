package services

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func searchArticle(aa []*models.Article, categorySlug, slug string) *models.Article {
	for _, a := range aa {
		if a.CategorySlug == categorySlug && a.Slug == slug {
			return a
		}
	}
	return nil
}

func ReadAllArticles() {
	infras.ReadAllArticles()
	infras.WaitIO()
}

func GetArticles() []*models.Article {
	infras.WaitIO()
    return infras.GetCachedArticles()
}

func FindArticle(categorySlug, slug string) *models.Article {
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	return searchArticle(aa, categorySlug, slug)
}

func AddArticle(a *models.Article) error {
	infras.WaitIO()
	ch := make(chan error)
	go infras.WriteArticle(a, ch)
	err := <-ch
	if err != nil {
		return err
	}
	infras.AwaitNextChange()
	return nil
}

func RemoveArticle(a *models.Article) error {
	infras.WaitIO()
	ch := make(chan error)
	infras.RemoveArticle(a, ch)
	err := <-ch
	if err != nil {
		return err
	}
	infras.AwaitNextChange()
	return nil
}

func IdentifyArticle(a *models.Article) *models.Article {
	aa := infras.GetCachedArticles()
	return searchArticle(aa, a.CategorySlug, a.Slug)
}

func ReplaceArticle(oldA, newA *models.Article) error {
	ch := make(chan error)
	infras.ReplaceArticle(oldA, newA, ch)
	err := <-ch
	if err != nil {
		return err
	}
	infras.AwaitNextChange()
	return nil
}
