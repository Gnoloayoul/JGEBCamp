package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateBYId(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime, art.Utime = now, now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateBYId(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})
	// 要不要检查是不是真的更新了？
	// 更新行数
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法 Id %d, author_id %d", art.Id, art.AuthorId)
	}
	return res.Error
}

func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 先操作制作库(此时应该是表)，后操作线上库(此时应该是表)
	// 在事务内部，这里采用了闭包形态
	// GORM 帮助我们管理了事务的生命周期
	// Begin， Rollback 和 Commit 都不需要操心
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var (
			id = art.Id
			err error
		)
		txDAO := NewGORMArticleDAO(tx)
		if  id > 0 {
			err = txDAO.UpdateBYId(ctx, art)

		}
	})
}

// Article
// [制作库]
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 限定长度：1024
	Title string `gorm:"type=varchar(1024)"`
	// BLOB：mysql 中适合存大文本数据的数据类型
	Content string `gorm:"type=BLOB"`
	// 仅仅给 AuthorId 上索引
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}
