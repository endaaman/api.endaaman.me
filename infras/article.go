package infras

import (
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/repositories"
)

type ArticleRepository struct {
}

func NewArticleRepository() repositories.ArticleRepository {
    repo := ArticleRepository{}
    return &repo
}

func (repo *ArticleRepository) FindAll() (aa []*models.Article, err error) {
    return nil, nil
}

func (repo *ArticleRepository) Find(word string) (aa []*models.Article, err error) {
    return nil, nil
}

func (repo *ArticleRepository) Create(a *models.Article) (*models.Article, error) {
    return nil, nil
}

func (repo *ArticleRepository) Update(a *models.Article) (*models.Article, error) {
    return nil, nil
}
