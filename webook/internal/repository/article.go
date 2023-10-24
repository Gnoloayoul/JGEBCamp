package repository

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)

}

type CachedArticleRepository struct {

}

func NewArticleRepository() ArticleRepository{
	return &CachedArticleRepository{

	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return 1, nil
}

