package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
)

//go:generate mockgen -source=./article_author.go -package=artrepomocks -destination=mocks/article_author_mock.go ArticleAuthorRepository
type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
