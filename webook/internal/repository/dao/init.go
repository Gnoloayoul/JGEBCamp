package dao

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{},
		&article.Article{},
		&article.PublishedArticle{},
		&Job{},
	)
}
