package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Article, error) {
	var res []Article
	err := dao.db.WithContext(ctx).
		Where("utime<?", start.UnixMilli()).
		Order("utime DESC").Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (dao *GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	var arts []Article
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("author_id = ?", author).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&arts).Error
	return arts, err
}

func (dao *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pub PublishedArticle
	err := dao.db.WithContext(ctx).
		Where("id = ?", id).
		First(&pub).Error
	return pub, err
}

func (dao *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, id int64, author int64, status uint8) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id=?", id, author).
			Update("status", status)
		if res.Error != nil {
			// 数据库有问题
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 要么 ID 是错的，要么作者不对
			// 后者情况下，就要小心，可能有人在搞系统
			// 没必要再用 ID 搜索数据库来区分这两种情况
			// 用 prometheus 打点，只要频繁出现，就告警，然后手工介入排查
			return ErrPossibleIncorrectAuthor
		}

		res = tx.Model(&PublishedArticle{}).
			Where("id=? AND author_id=?", id, author).
			Update("status", status)
		if res.Error != nil {
			// 数据库有问题
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})
}

func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	tx := dao.db.WithContext(ctx).Begin()
	now := time.Now().UnixMilli()
	defer tx.Rollback()
	txDAO := NewGORMArticleDAO(tx)

	// 先操作制作库(此时应该是表)，后操作线上库(此时应该是表)
	// 在事务内部，这里采用了闭包形态
	// GORM 帮助我们管理了事务的生命周期
	// Begin， Rollback 和 Commit 都不需要操心
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = txDAO.Insert(ctx, art)
	} else {
		err = txDAO.UpdateById(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	publishedArt := PublishedArticle(art)
	publishedArt.Utime = now
	publishedArt.Ctime = now
	err = tx.Clauses(clause.OnConflict{

		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		}),
	}).Create(&publishedArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, tx.Error
}

func (dao *GORMArticleDAO) SyncClosure(ctx context.Context, art Article) (int64, error) {
	var (
		id = art.Id
	)
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		txDAO := NewGORMArticleDAO(tx)
		if id == 0 {
			id, err = txDAO.Insert(ctx, art)
		} else {
			err = txDAO.UpdateById(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		publishedArt := art
		publishedArt.Utime = now
		publishedArt.Ctime = now
		return tx.Clauses(clause.OnConflict{

			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   art.Title,
				"content": art.Content,
				"status":  art.Status,
				"utime":   now,
			}),
		}).Create(&publishedArt).Error
	})
	return id, err
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime, art.Utime = now, now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		})
	// 要不要检查是不是真的更新了？
	// 更新行数
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
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
			"status":  art.Status,
			"utime":   now,
		}),
	}).Create(&art).Error
	// 在 Mysql 里最终生成的语句是这
	// INSERT xxx ON DUPLICATE KEY UPDATE XXX
	// 正常而言，一条 SQL， 是不需要开启的
	// 但要小心 auto commit： 自动提交
	return err
}
