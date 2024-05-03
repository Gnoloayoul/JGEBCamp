package week11

import (
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestConnPoolV1(t *testing.T) {
	d1, err := gorm.Open(mysql.Open("root:root@tcp(124.156.150.17:13316)/testBase"))
	require.NoError(t, err)
	err = d1.AutoMigrate(&TestTable{})
	require.NoError(t, err)

	d2, err := gorm.Open(mysql.Open("root:root@tcp(124.156.150.17:13316)/testBase_NewOne"))
	require.NoError(t, err)
	err = d2.AutoMigrate(&TestTable{})
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWriterPool{
			src:     d1.ConnPool,
			dst:     d2.ConnPool,
			pattern: atomicx.NewValueOf(patternDstFirst),
		},
	}))
	require.NoError(t, err)

	t.Log(db)

	err = db.Create(&TestTable{
		Biz:   "DWtest",
		BizId: 112233,
	}).Error
	require.NoError(t, err)

	err = db.Transaction(func(tx *gorm.DB) error {
		err1 := tx.Create(&TestTable{
			BizId: 1122334,
			Biz:   "DWtest_tx",
		}).Error

		return err1
	})
	require.NoError(t, err)

	err = db.Model(&TestTable{}).Where("id > ?", 0).Updates(map[string]any{
		"biz_id": 789,
	}).Error
	require.NoError(t, err)
}

type TestTable struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}
