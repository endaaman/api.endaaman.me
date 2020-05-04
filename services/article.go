package services

import (
	"fmt"

	"github.com/endaaman/api.endaaman.me/infras"
	"github.com/endaaman/api.endaaman.me/models"
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

func GetArticles(includePrivate bool) []*models.Article {
	infras.WaitIO()
	cached := infras.GetCachedArticles()
	if includePrivate {
		return cached
	}
	aa := make([]*models.Article, 0)
	for _, a := range cached {
		if !a.Private {
			aa = append(aa, a)
		}
	}
	return aa
}

func FindArticle(categorySlug, slug string) *models.Article {
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	return searchArticle(aa, categorySlug, slug)
}

func AddArticle(a *models.Article) error {
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	oldA := searchArticle(aa, a.CategorySlug, a.Slug)
	if oldA != nil {
		return fmt.Errorf("Article `%s/%s` already exists.", a.CategorySlug, a.Slug)
	}
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
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	return searchArticle(aa, a.CategorySlug, a.Slug)
}

func UpdateArticle(oldA, newA *models.Article) error {
	infras.WaitIO()
	aa := infras.GetCachedArticles()
	slugChanged := newA.JointedSlug() != oldA.JointedSlug()
	if slugChanged {
		existingA := searchArticle(aa, newA.CategorySlug, newA.Slug)
		if existingA != nil {
			return fmt.Errorf("Article `%s` already exists.", newA.JointedSlug())
		}
	}
	ch := make(chan error)
	infras.UpdateArticle(oldA, newA, ch)
	err := <-ch
	if err != nil {
		return err
	}
	infras.AwaitNextChange()
	return nil
}
