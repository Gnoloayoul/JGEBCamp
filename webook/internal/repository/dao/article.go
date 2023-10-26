package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, article Article) error
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
	err := dao.db.WithContext(ctx).Model(&art).
		Where("id=?", art.Id).
		Updates(map[string]any{
			"title": art.Title,
			"content": art.Content,
			"utime": art.Utime,
	})
	return err
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
	Ctime int64
	Utime int64
}