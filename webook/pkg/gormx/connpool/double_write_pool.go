package connpool

import (
	"context"
	"database/sql"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写模式")

type DoubleWriterPool struct {
	src gorm.ConnPool
	dst gorm.ConnPool
	pattern *atomicx.Value[string]
}

func (d *DoubleWriterPool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case patternSrcOnly:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriterPoolTx{
			src: tx,
			pattern: pattern,
		}, err
	case patternSrcFirst:
		srcTx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dstTx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// 记录日志，然后不做处理

			// 可以考虑回滚
			// err = srcTx.Rollback()
			// return err
		}
		return &DoubleWriterPoolTx{
			src: srcTx,
			dst: dstTx,
			pattern: pattern,
		}, nil
	case patternDstOnly:
		tx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriterPoolTx{
			src: tx,
			pattern: pattern,
		}, err
	case patternDstFirst:
		dstTx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		srcTx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// 记录日志，然后不做处理

			// 可以考虑回滚
			// err = srcTx.Rollback()
			// return err
		}
		return &DoubleWriterPoolTx{
			src: srcTx,
			dst: dstTx,
			pattern: pattern,
		}, nil
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriterPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (d *DoubleWriterPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// src写失败，不被认为是失败
		}
		return res, err
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPool) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
	// 能不能再有事务未提交的情况加，我禁止修改?
	// 能，需要维持一个已开事务的计数，要用锁，因此有严重的性能问题
}

type DoubleWriterPoolTx struct {
	src *sql.Tx
	dst *sql.Tx
	pattern string
}

func (d *DoubleWriterPoolTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (d *DoubleWriterPoolTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.dst == nil {
			return res, err
		}
		_, err = d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.src == nil {
			return res, err
		}
		_, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// src写失败，不被认为是失败
		}
		return res, err
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case patternSrcOnly, patternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
	}
}

func (d *DoubleWriterPoolTx) Commit() error {
	switch d.pattern {
	case patternSrcOnly:
		return d.src.Commit()
	case patternSrcFirst:
		err := d.src.Commit()
		if err != nil {
			return err
		}
		if d.dst != nil {
			err = d.dst.Commit()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	case patternDstOnly:
		return d.dst.Commit()
	case patternDstFirst:
		err := d.dst.Commit()
		if err != nil {
			return err
		}
		if d.src != nil {
			err = d.src.Commit()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriterPoolTx) Rollback() error {
	switch d.pattern {
	case patternSrcOnly:
		return d.src.Rollback()
	case patternSrcFirst:
		err := d.src.Rollback()
		if err != nil {
			return err
		}
		if d.dst != nil {
			err = d.dst.Rollback()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	case patternDstOnly:
		return d.dst.Rollback()
	case patternDstFirst:
		err := d.dst.Rollback()
		if err != nil {
			return err
		}
		if d.src != nil {
			err = d.src.Rollback()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	default:
		return errUnknownPattern
	}
}

const (
	patternDstOnly  = "DST_ONLY"
	patternSrcOnly  = "SRC_ONLY"
	patternDstFirst = "DST_FIRST"
	patternSrcFirst = "SRC_FIRST"
)