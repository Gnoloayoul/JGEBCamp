package client

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	"google.golang.org/grpc"
)

type InteractiveRepositoryAdapter struct {
	repo repository.InteractiveRepository
}

func (i *InteractiveRepositoryAdapter) IncrReadCnt(ctx context.Context, in *intrRepov1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrRepov1.IncrReadCntResponse, error) {
	err := i.repo.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrRepov1.IncrReadCntResponse{}, err
}

func (i *InteractiveRepositoryAdapter) BatchIncrReadCnt(ctx context.Context, in *intrRepov1.BatchIncrReadCntRequest, opts ...grpc.CallOption) (*intrRepov1.BatchIncrReadCntResponse, error) {
	err := i.repo.BatchIncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrRepov1.BatchIncrReadCntResponse{}, err
}

func (i *InteractiveRepositoryAdapter) IncrLike(ctx context.Context, in *intrRepov1.IncrLikeRequest, opts ...grpc.CallOption) (*intrRepov1.IncrLikeResponse, error) {
	err := i.repo.IncrLike(ctx, in.GetBiz(), in.GetBizId(), in.GetBizId())
	return &intrRepov1.IncrLikeResponse{}, err
}

func (i *InteractiveRepositoryAdapter) DecrLike(ctx context.Context, in *intrRepov1.DecrLikeRequest, opts ...grpc.CallOption) (*intrRepov1.DecrLikeResponse, error) {
	err := i.repo.DecrLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrRepov1.DecrLikeResponse{}, err
}

func (i *InteractiveRepositoryAdapter) AddCollectionItem(ctx context.Context, in *intrRepov1.AddCollectionItemRequest, opts ...grpc.CallOption) (*intrRepov1.AddCollectionItemResponse, error) {
	err := i.repo.AddCollectionItem(ctx, in.GetBiz(), in.GetBizId(), in.GetCid(), in.GetUid())
	return &intrRepov1.AddCollectionItemResponse{}, err
}

func (i *InteractiveRepositoryAdapter) Get(ctx context.Context, in *intrRepov1.GetRequest, opts ...grpc.CallOption) (*intrRepov1.GetResponse, error) {
	intr, err := i.repo.Get(ctx, in.GetBiz(), in.GetBizId())
	return &intrRepov1.GetResponse{
		Intr: i.toDTO(intr),
	}, err
}

func (i *InteractiveRepositoryAdapter) Liked(ctx context.Context, in *intrRepov1.LikedRequest, opts ...grpc.CallOption) (*intrRepov1.LikedResponse, error) {
	_, err := i.repo.Liked(ctx, in.GetBiz(), in.GetId(), in.GetUid())
	return &intrRepov1.LikedResponse{}, err
}

func (i *InteractiveRepositoryAdapter) Collected(ctx context.Context, in *intrRepov1.CollectedRequest, opts ...grpc.CallOption) (*intrRepov1.CollectedResponse, error) {
	_, err := i.repo.Collected(ctx, in.GetBiz(), in.GetId(), in.GetUid())
	return &intrRepov1.CollectedResponse{}, err
}

func (i *InteractiveRepositoryAdapter) GetByIds(ctx context.Context, in *intrRepov1.GetByIdsRequest, opts ...grpc.CallOption) (*intrRepov1.GetByIdsResponse, error) {
	intrs, err := i.repo.GetByIds(ctx, in.GetBiz(), in.GetIds())
	res := make([]*intrRepov1.Interactive, len(intrs))
	for k, v := range intrs {
		res[k] = i.toDTO(v)
	}
	return &intrRepov1.GetByIdsResponse{
		Intrs: res,
	}, err
}

func NewInteractiveRepositoryAdapter(repo repository.InteractiveRepository) *InteractiveRepositoryAdapter{
	return &InteractiveRepositoryAdapter{
		repo: repo,
	}
}

func (i *InteractiveRepositoryAdapter) toDTO(intr domain.Interactive) *intrRepov1.Interactive {
	return &intrRepov1.Interactive{
		Biz: intr.Biz,
		BizId: intr.BizId,
		ReadCnt: intr.ReadCnt,
		LikeCnt: intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked: intr.Liked,
		Collected: intr.Collected,
	}
}