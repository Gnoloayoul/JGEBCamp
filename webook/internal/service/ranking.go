package service

import (
	"context"
	"errors"
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"time"
)

//go:generate mockgen -source=./ranking.go -package=svcmocks -destination=mocks/ranking_mock.go RankingService
type RankingService interface {
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	artSvc    ArticleService
	intrSvc   intrv1.InteractiveServiceClient
	repo      repository.RankingRepository
	batchSize int
	n         int
	// scoreFunc 不能返回负数
	// 实际的算法公式
	scoreFunc func(t time.Time, likeCnt int64) float64
	load      int64
}

func NewBatchRankingService(artSvc ArticleService,
	repo repository.RankingRepository,
	intrSvc intrv1.InteractiveServiceClient) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		repo:      repo,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(float64(sec+2), 1.5)
		},
	}
}

func (svc *BatchRankingService) TopN(ctx context.Context) error {
	// 调下面的来topN
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	// 在这里，存起来
	return svc.repo.ReplaceTopN(ctx, arts)
}

func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	// 我只取七天内的数据
	now := time.Now()
	// 先拿一批数据
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})

	for {
		// 先拿一批
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})
		// 要去找到对应的点赞数据
		intrs, err := svc.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article", BizIds: ids,
		})
		if err != nil {
			return nil, err
		}
		if len(intrs.Intrs) == 0 {
			return nil, errors.New("没有数据")
		}
		// 合并计算 score
		// 排序
		for _, art := range arts {
			intr := intrs.Intrs[art.Id]
			score := svc.scoreFunc(art.Utime, intr.LikeCnt)
			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})
			// 这种写法，要求 topN 已经满了
			if err == queue.ErrOutOfCapacity {
				val, _ := topN.Dequeue()
				if val.score < score {
					_ = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}
		// 一批已经处理完了，问题来了，我要不要进入下一批？我怎么知道还有没有？
		if len(arts) < svc.batchSize ||
			now.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			// 我这一批都没取够，我当然可以肯定没有下一批了
			break
		}
		// 这边要更新 offset
		offset = offset + len(arts)
	}
	// 最后得出结果
	res := make([]domain.Article, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		res[i] = val.art
	}
	return res, nil
}
