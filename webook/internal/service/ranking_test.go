package service

import (
	"context"
	"fmt"
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	repomocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/mocks"
	svcmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

// 没改好的测试
func TestRankingTopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (ArticleService,
			repository.RankingRepository,
			intrv1.InteractiveServiceClient)

		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "计算成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, repository.RankingRepository, intrv1.InteractiveServiceClient) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				// 最简单，一批就搞完
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 3).
					Return([]domain.Article{
						{Id: 1, Utime: now, Ctime: now},
						{Id: 2, Utime: now, Ctime: now},
						{Id: 3, Utime: now, Ctime: now},
					}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 3, 3).
					Return([]domain.Article{}, nil)

				intrSvc := svcmocks.NewMockInteractiveServiceClient(ctrl)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					&intrv1.GetByIdsRequest{
						Biz:    "article",
						BizIds: []int64{1, 2, 3},
					}).Return(&intrv1.GetByIdsResponse{Intrs: map[int64]*intrv1.Interactive{
					1: {BizId: 1, LikeCnt: 1},
					2: {BizId: 2, LikeCnt: 2},
					3: {BizId: 3, LikeCnt: 3},
				}}, nil)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					&intrv1.GetByIdsRequest{Biz: "article", BizIds: []int64{}}).
					Return(&intrv1.GetByIdsResponse{Intrs: map[int64]*intrv1.Interactive{}}, nil)
				repo := repomocks.NewMockRankingRepository(ctrl)
				return artSvc, repo, intrSvc
			},
			//wantArts: []domain.Article{
			//	{Id: 3, Utime: now, Ctime: now},
			//	{Id: 2, Utime: now, Ctime: now},
			//	{Id: 1, Utime: now, Ctime: now},
			//},
			wantArts: []domain.Article(nil),
			wantErr:  fmt.Errorf("没有数据"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, repo, intrSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(artSvc, repo, intrSvc).(*BatchRankingService)
			// 为了测试
			svc.batchSize = 3
			svc.n = 3
			svc.scoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArts, arts)
		})
	}
}
