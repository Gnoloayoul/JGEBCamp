package article

import (
	"context"
	"gorm.io/gorm"
)

type AuthorDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateBYId(ctx context.Context, article Article) error
}

func NewAuthorDAO(db *gorm.DB) AuthorDao {
	panic("v")
}
