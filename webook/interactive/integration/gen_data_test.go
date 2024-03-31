package integration

import (
	_ "embed"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/integration/startup"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

//go:embed init.sql
var initSQL string

func TestGenSQL(t *testing.T) {
	file, err := os.OpenFile("data.sql",
		os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(initSQL)
	require.NoError(t, err)

	const prefix = "INSERT INTO `interactives`(`biz_id`, `biz`, `read_cnt`, `collect_cnt`, `like_cnt`, `ctime`, `utime`)\nVALUES"
	const rowNum = 1000

	now := time.Now().UnixMilli()
	_, err = file.WriteString(prefix)

	for i := 0; i < rowNum; i++ {
		if i > 0 {
			file.Write([]byte{',', '\n'})
		}

		file.Write([]byte{'('})

		// biz_id
		file.WriteString(strconv.Itoa(i + 1))

		// biz: "test"
		file.WriteString(`,"test",`)

		// read_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// collect_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// like_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// ctime
		file.WriteString(strconv.FormatInt(now, 10))
		file.Write([]byte{','})

		// utime
		file.WriteString(strconv.FormatInt(now, 10))

		file.Write([]byte{')'})
	}
}

func TestGenData(t *testing.T) {
	db := startup.InitTestDB()
	for i := 0; i < 10; i++ {
		const bitchSize = 100
		data := make([]dao.Interactive, 0, bitchSize)
		now := time.Now().UnixMilli()
		for j := 0; j < bitchSize; j++ {
			data = append(data, dao.Interactive{
				BizId: int64(i * bitchSize + j + 1),
				Biz: "test",
				ReadCnt: rand.Int63(),
				CollectCnt: rand.Int63(),
				LikeCnt: rand.Int63(),
				Ctime: now,
				Utime: now,
			})
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			err := tx.Create(data).Error
			require.NoError(t, err)
			return err
		})
		require.NoError(t, err)
	}

}