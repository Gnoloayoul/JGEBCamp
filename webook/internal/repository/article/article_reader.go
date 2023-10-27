package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
)

type ArticleReaderRepository interface {
	// 有就更新，没有就新建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
