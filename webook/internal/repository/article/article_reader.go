package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
)

//go:generate mockgen -source=./article_reader.go -package=artrepomocks -destination=mocks/article_reader_mock.go ArticleReaderRepository
type ArticleReaderRepository interface {
	// 有就更新，没有就新建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
