package article

import (
	"context"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./author_dao.go -package=articlemocks -destination=mocks/author_dao_mock.go ArticleDAO
type AuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

func NewAuthorDAO(db *gorm.DB) AuthorDAO {
	panic("implement me")
}
