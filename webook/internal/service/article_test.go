package service

import (
	"context"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article"
	artrepomocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article/mocks"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)

		art domain.Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)

				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return author, reader
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)

				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)

				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "在制作库就失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)

				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存失败"))

				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("保存失败"),
		},
		{
			name: "保存至制作库成功，但重试到线上库成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)

				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("保存失败"))

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)

				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "保存至制作库成功，但重试全失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)

				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Times(3).Return(int64(0), errors.New("保存失败"))

				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("保存失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleService(author, reader, &logger.NopLogger{})
			id, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
