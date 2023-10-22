package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// ArticleTestSuite
// 集成测试套件
type ArticleTestSuite struct {
	suite.Suite
}


func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello， 这里是测试套件")
}


func TestAriticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}