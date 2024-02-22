package integration

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/integration/startup"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"time"
)

type InteractiveTestSuite struct {
	suite.Suite
	db *gorm.DB
	rdb redis.Cmdable
}

func (s *InteractiveTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
}

func (s *InteractiveTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()

	err := s.db.Exec("TRUNCATE TABLE 'interactive'").Error
	assert.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE 'user_like_bizs'").Error
	assert.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE 'user_collection_bizs'").Error
	assert.NoError(s.T(), err)

	// clear Redis
	err = s.rdb.FlushDB(ctx).Err()
	assert.NoError(s.T(), err)
}