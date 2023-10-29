package article

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateBYId(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, art PublishedArticle) error
	SyncStatus(ctx context.Context, id int64, author int64, status uint8) error
}


type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, id int64, author int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id=?", id, author).
			Updates(map[string]any{
				"status": status,
				"utime": now,
		})
		if res.Error != nil {
			// 数据库有问题
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 要么 ID 是错的，要么作者不对
			// 后者情况下，就要小心，可能有人在搞系统
			// 没必要再用 ID 搜索数据库来区分这两种情况
			// 用 prometheus 打点，只要频繁出现，就告警，然后手工介入排查
			return fmt.Errorf("可能有人在搞系统, 误操作非自己的文章， uid：%d, aid: %d", author, id)
		}
		return tx.Model(&Article{}).
			Where("id=?", id).
			Updates(map[string]any{
				"status": status,
				"utime": now,
			}).Error
	})
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
			"status": art.Status,
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
	var (
		id = art.Id
	)
	// tx -> Transaction, 也有人缩写成 trx
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDAO(tx)
		if  id > 0 {
			err = txDAO.UpdateBYId(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}

		// 要操作线上库了
		return txDAO.Upsert(ctx, PublishedArticle{Article: art})
	})
	return id, err
}

// Upsert
// Insert or Update
// 在 db 上实现
func (dao *GORMArticleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime, art.Utime = now, now
	// 插入
	// OnConflict 数据冲突了
	err := dao.db.Clauses(clause.OnConflict{
		// 用 GORM—Mysql 只需要关心这里
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"status": art.Status,
			"utime":   now,
	}),
	}).Create(&art).Error
	// 在 Mysql 里最终生成的语句是这
	// INSERT xxx ON DUPLICATE KEY UPDATE XXX
	// 正常而言，一条 SQL， 是不需要开启的
	// 但要小心 auto commit： 自动提交
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
	Status uint8
	Ctime    int64
	Utime    int64
}
