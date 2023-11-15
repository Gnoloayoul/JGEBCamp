package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ArticleGORMFuncTestSuite struct {
	suite.Suite
	server *gin.Engine
	db *gorm.DB
}

func (s *ArticleGORMFuncTestSuite) SetupSuite() {

}

func (s *ArticleGORMFuncTestSuite) TearDownTest() {

}

func (s *ArticleGORMFuncTestSuite)Test_findN() {

}

func TestRun(t *testing.T) {
	suite.Run(t, new(ArticleGORMFuncTestSuite))
}