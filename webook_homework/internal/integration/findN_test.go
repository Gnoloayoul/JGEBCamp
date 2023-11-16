package integration

import (
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/internal/integration/startup"
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/internal/repository/dao/article"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook_homework/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ArticleGORMFuncTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleGORMFuncTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(context *gin.Context) {
		context.Set("claims", &ijwt.UserClaims{
			Id: 123,
		})
		context.Next()
	})
	s.db = startup.InitTestDB()
	hdl := startup.InitArticleHandler(article.NewGORMArticleDAO(s.db))
	hdl.RegisterRoutes(s.server)
}

func (s *ArticleGORMFuncTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `articles`").Error
	assert.NoError(s.T(), err)
	s.db.Exec("TRUNCATE TABLE `published_articles`")
}

func (s *ArticleGORMFuncTestSuite) Test_findN() {
	t := s.T()
	testCase := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req
	}
}

func TestRun(t *testing.T) {
	suite.Run(t, new(ArticleGORMFuncTestSuite))
}
