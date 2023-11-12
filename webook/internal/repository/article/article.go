package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	dao "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

// repository 应该是 cache 与 dao （或者一些高级操作的）的胶水
// 事务概念应该在 DAO 这一层

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	//FindById(ctx context.Context, id int64) domain.Article
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) (int64, error)
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO

	// V1 操作两个 dao
	readerDao dao.ReaderDao
	authorDao dao.AuthorDAO

	// V2用
	// 耦合了 DAO 操作的东西
	// 正常情况下，如果要在 repository 上操作事务
	// 要么只能利用 db 开始事务之后，创建基于事务的 dao
	// 或者，直接去掉 DAO 这一层，在 repository 的实现中，直接操作数据库 db
	db *gorm.DB

	cache cache.ArticleCache
	l     logger.LoggerV1
}

func (c *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(res), nil
}
