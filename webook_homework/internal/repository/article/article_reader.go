package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/internal/domain"
)

type ArticleReaderRepository interface {
	// Save 有就更新，没有就新建，即 upsert 的语义
	Save(ctx context.Context, art domain.Article) (int64, error)
}
