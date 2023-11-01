package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	dao "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/article"
	"gorm.io/gorm"
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
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO

	// V1 操作两个 dao
	readerDao dao.ReaderDao
	authorDao dao.AuthorDao

	// V2用
	// 耦合了 DAO 操作的东西
	// 正常情况下，如果要在 repository 上操作事务
	// 要么只能利用 db 开始事务之后，创建基于事务的 dao
	// 或者，直接去掉 DAO 这一层，在 repository 的实现中，直接操作数据库 db
	db *gorm.DB
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateBYId(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) (int64, error) {
	return c.dao.SyncStatus(ctx, id, author, status.ToUint8())
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = c.authorDao.UpdateBYId(ctx, artn)
	} else {
		id, err = c.authorDao.Insert(ctx, artn)
	}
	if err != nil {
		return id, nil
	}
	// 操作线上库，保存数据，同步过来
	// 考虑到线上库可能有，也可能没有，因此需要 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么更新，不然就是插入
	err = c.readerDao.Upsert(ctx, artn)
	return id, err
}

// SyncV2
// 尝试在 repository 层面上解决事务问题
// 确保保存到线上库和制作库同时成功或者同时失败
func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	// 开始一个事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer tx.Rollback()
	// 利用 tx 来构建 DAO
	author := dao.NewAuthorDAO(tx)
	reader := dao.NewReaderDAO(tx)

	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = author.UpdateBYId(ctx, artn)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {
		return id, nil
	}
	// 操作线上库，保存数据，同步过来
	// 考虑到线上库可能有，也可能没有，因此需要 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么更新，不然就是插入
	err = reader.Upsert(ctx, artn)
	// 执行成功， 直接提交
	tx.Commit()
	return id, err
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
