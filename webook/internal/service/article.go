package service

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	WithDraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo article.ArticleRepository

	// v1用 与上面的互斥
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger.LoggerV1
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) WithDraw(ctx context.Context, art domain.Article) error {
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	// 重试机制
	// TODO: 将重试改为装饰器模式
	for i := 0; i < 3; i++ {
		// 当数据库发生问题，先等待一下，看数据库能不能恢复过来
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败， 保存到制作库",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败， 重试侧底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		// 这时该接入告警系统，手工处理一下
		// 走异步，直接保存到本地
		// 走 Canle
		// 打 MQ
	}
	return id, err
}
