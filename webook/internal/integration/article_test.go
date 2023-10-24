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

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct{
		name string

		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after func(t *testing.T)

		art Article

		wantCode int
		wantRes Result[int64]

	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			// 构造请求
			// 执行
			// 验证结果

		})
	}
}


func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello， 这里是测试套件")
}


func TestAriticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Title string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data T `json:data`
}