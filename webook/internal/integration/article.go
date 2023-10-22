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

func TestAriticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}