package repositories

import (
	"github.com/endaaman/api.endaaman.me/models"
)

type ArticleRepository interface {
    FindAll() (aa []*models.Article, err error)
    Find(word string) (aa []*models.Article, err error)
    Create(a *models.Article) (*models.Article, error)
    Update(a *models.Article) (*models.Article, error)
}
