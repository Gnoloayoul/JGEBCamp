package article

import (
	"context"
)

type AuthorDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateBYId(ctx context.Context, article Article) error
}