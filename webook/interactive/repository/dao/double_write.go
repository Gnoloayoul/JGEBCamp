package dao

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	patternDstOnly  = "DST_ONLY"
	patternSrcOnly  = "SRC_ONLY"
	patternDstFirst = "DST_FIRST"
	patternSrcFirst = "SRC_FIRST"
)

type DoubleWriterDAO struct {
	src     InteractiveDAO
	dst     InteractiveDAO
	pattern *atomicx.Value[string]
}

func NewDoubleWriterDAOV1(src, dst *gorm.DB) *DoubleWriterDAO {
	return &DoubleWriterDAO{
		src:     NewGORMInteractiveDAO(src),
		dst:     NewGORMInteractiveDAO(dst),
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func NewDoubleWriterDAO(src, dst InteractiveDAO) *DoubleWriterDAO {
	return &DoubleWriterDAO{
		src:     src,
		dst:     dst,
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func (d *DoubleWriterDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}

func (d *DoubleWriterDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.IncrReadCnt(ctx, biz, bizId)
	case patternDstOnly:
		return d.dst.IncrReadCnt(ctx, biz, bizId)
	case patternSrcFirst:
		err := d.src.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}

		err = d.dst.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			// 记日志
		}

		return nil
	case patternDstFirst:
		err := d.dst.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}

		err = d.src.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			// 记日志
		}

		return err
	default:
		return errors.New("未知的双写模式")
	}
}

func (d *DoubleWriterDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	panic("implement me")
}

func (d *DoubleWriterDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	panic("implement me")
}

func (d *DoubleWriterDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	panic("implement me")
}

func (d *DoubleWriterDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	panic("implement me")
}

func (d *DoubleWriterDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	panic("implement me")
}

func (d *DoubleWriterDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	panic("implement me")
}

func (d *DoubleWriterDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	panic("implement me")
}

func (d *DoubleWriterDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	panic("implement me")
}
