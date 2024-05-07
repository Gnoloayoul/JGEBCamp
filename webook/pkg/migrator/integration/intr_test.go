package integration

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/integration/startup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

type InteractiveTestSuite struct {
	suite.Suite
	srcDB *gorm.DB
	intrDB *gorm.DB
}

func (i *InteractiveTestSuite) SetupSuite() {
	i.srcDB = startup.InitSrcDB()
	err := i.srcDB.AutoMigrate(&Interactive{})
	assert.NoError(i.T(), err)
	i.intrDB = startup.InitDstDB()
	err = i.intrDB.AutoMigrate(&Interactive{})
	assert.NoError(i.T(), err)
}

func (i *InteractiveTestSuite) TearDownTest() {
	i.srcDB.Exec("TRUNCATE TABLE interactives")
	i.intrDB.Exec("TRUNCATE TABLE interactives")
}

func (i *InteractiveTestSuite) TestValidator() {
	testCases := []struct{
		name string
		before func(t *testing.T)
		after func(t *testing.T)
		mock func(ctrl *gomock.Controller) events2.Producer

		wantErr error
	}{
		{},

	}
	for _, tc :=

}

func TestInteractive(t *testing.T) {
	suite.Run(t, &InteractiveTestSuite{})
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Ctime      int64
	Utime      int64
}

func (i Interactive) ID() int64 {
	return i.Id
}

func (i Interactive) TableName() string {
	return "interactives"
}

func (i Interactive) CompareTo(entity migrator.Entity) bool {
	dst := entity.(migrator.Entity)
	return i == dst
}