package repository

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
)

//go:generate mockgen -source=./interactive.go -package=repomocks -destination=mocks/interactive_mock.go InteractiveRepository
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context,
		biz []string, bizId []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.LoggerV1
}


func (c *CachedReadCntRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	vals, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Interactive, domain.Interactive](vals,
		func(idx int, src dao.Interactive) domain.Interactive {
			return c.toDomain(src)
		}), nil

}

func (c *CachedReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrDataNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrDataNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	// 先插入点赞，然后更新点赞计数，更新缓存
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	// 这种做法，你需要在 repository 层面上维持住事务
	//c.dao.IncrLikeCnt()
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// 要考虑缓存方案了
	// 这两个操作能不能换顺序？ —— 不能
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//go func() {
	//	c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
	//}()
	//return err

	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

// BatchIncrReadCnt
// bizs 和 ids 的长度必须相等
func (c *CachedReadCntRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, bizId []int64) error {
	// 我在这里要不要检测 bizs 和 ids 的长度是否相等？
	err := c.dao.BatchIncrReadCnt(ctx, bizs, bizId)
	if err != nil {
		return err
	}
	// 你也要批量的去修改 redis，所以就要去改 lua 脚本
	// c.cache.IncrReadCntIfPresent()
	// TODO, 等我写新的 lua 脚本/或者用 pipeline
	return nil
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	// 这个地方，你要不要考虑缓存收藏夹？
	// 以及收藏夹里面的内容
	// 用户会频繁访问他的收藏夹，那么你就应该缓存，不然你就不需要
	// 一个东西要不要缓存，你就看用户会不会频繁访问（反复访问）
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Cid:   cid,
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	// 收藏个数（有多少个人收藏了这个 biz + bizId)
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	// 要从缓存拿出来阅读数，点赞数和收藏数
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}

	ie, err := c.dao.Get(ctx, biz, bizId)
	if err == nil  {
		res := c.toDomain(ie)
		if er := c.cache.Set(ctx, biz, bizId, res); er != nil {
			c.l.Error("回写缓存失败",
				logger.Int64("bizId", bizId),
				logger.String("biz", biz),
				logger.Error(er))
		}
		return res, nil
	}
	return domain.Interactive{}, err
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz: intr.Biz,
		BizId:      intr.BizId,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO,
	cache cache.InteractiveCache, l logger.LoggerV1) InteractiveRepository {
	return &CachedReadCntRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}


// 正常来说，参数必然不用指针：方法不要修改参数，通过返回值来修改参数
// 返回值就看情况。如果是指针实现了接口，那么就返回指针
// 如果返回值很大，你不想值传递引发复制问题，那么还是返回指针
// 返回结构体

// 最简原则：
// 1. 接收器永远用指针
// 2. 输入输出都用结构体

