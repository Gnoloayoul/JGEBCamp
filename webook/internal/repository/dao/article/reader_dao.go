package article

import(
	"context"
	"gorm.io/gorm"
)

type ReaderDao interface {
	Upsert(ctx context.Context, art Article) error
}

// PublishedArticle
// 代表线上库
// (同库不同表)
type PublishedArticle struct {
	Article
}

func NewReaderDAO(db *gorm.DB) ReaderDao {
	panic("v")
}