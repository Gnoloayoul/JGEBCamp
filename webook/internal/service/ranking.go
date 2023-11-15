package service

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article"
)

type RankingService interface {
	TopN(ctx context.Context) error
}

func (svc *BatchRankingService) TopN(ctx interface{}) error {
	// 调下面的来topN
}

func (svc *BatchRankingService) topN(ctx interface{}) ([]) {
	// 先拿一批
	for {
		svc.artSvc
	}
}

type BatchRankingService struct {
	artSvc article.ArticleService
	intrSvc
}



