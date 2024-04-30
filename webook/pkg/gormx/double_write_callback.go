package gormx

import (
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

type DoubleWriteCallback struct {
	src     *gorm.DB
	dst     *gorm.DB
	pattern *atomicx.Value[string]
}

func (d *DoubleWriteCallback) create() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// Todo: 使用 gorm callback 实现对两个 db 的操作（双写）
		// 难点：这里只能有一个 db 进来，不是 src 就是 dst
		// 而且也做不到动态切换
		//d.src.Create(db.Statement.Model).Error
	}
}
