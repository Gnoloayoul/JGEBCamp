package repository

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
)

//go:generate mockgen -source=./ranking.go -package=repomocks -destination=mocks/ranking_mock.go RankingRepository
type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	// 使用具体实现，可读性更好，对测试不友好，因为咩有面向接口编程
	redis *cache.RankingRedisCache
	local *cache.RankingLocalCache
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	_ = c.local.Set(ctx, arts)
	return c.redis.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	data, err := c.local.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		c.local.Set(ctx, data)
	} else {
		return c.local.ForceGet(ctx)
	}
	return data, err
}

func NewCachedRankingRepository(
	redis *cache.RankingRedisCache,
	local *cache.RankingLocalCache,
) RankingRepository {
	return &CachedRankingRepository{local: local, redis: redis}
}
