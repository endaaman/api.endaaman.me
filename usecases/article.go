package usecases

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/repositories"
)

type ArticleUsecase interface {
    Search(string) (a []*models.Article, err error)
    View() (a []*models.Article, err error)
    Add(*models.Article) (err error)
}

type articleUsecase struct {
    repo repositories.ArticleRepository
}

func NewArticleUsecase(repo repositories.ArticleRepository) ArticleUsecase {
    u := articleUsecase{repo: repo}
    return &u
}

func (usecase *articleUsecase) Search(word string) (aa []*models.Article, err error) {
	return nil, nil
}

func (usecase *articleUsecase) View() (aa []*models.Article, err error) {
	return nil, nil
}

func (usecase *articleUsecase) Add(a *models.Article) (err error) {
	return nil
}
