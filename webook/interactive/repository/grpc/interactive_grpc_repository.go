package grpc

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	"google.golang.org/grpc"
)

type InteractiveGRPCRepositoryServer struct {
	repo repository.InteractiveRepository
	intrRepov1.UnimplementedInteractiveRepositoryServer
}

func (i *InteractiveGRPCRepositoryServer) IncrReadCnt(ctx context.Context, request *intrRepov1.IncrReadCntRequest) (*intrRepov1.IncrReadCntResponse, error) {
	err := i.repo.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrRepov1.IncrReadCntResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) BatchIncrReadCnt(ctx context.Context, request *intrRepov1.BatchIncrReadCntRequest) (*intrRepov1.BatchIncrReadCntResponse, error) {
	err := i.repo.BatchIncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrRepov1.BatchIncrReadCntResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) IncrLike(ctx context.Context, request *intrRepov1.IncrLikeRequest) (*intrRepov1.IncrLikeResponse, error) {
	err := i.repo.IncrLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrRepov1.IncrLikeResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) DecrLike(ctx context.Context, request *intrRepov1.DecrLikeRequest) (*intrRepov1.DecrLikeResponse, error) {
	err := i.repo.DecrLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrRepov1.DecrLikeResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) AddCollectionItem(ctx context.Context, request *intrRepov1.AddCollectionItemRequest) (*intrRepov1.AddCollectionItemResponse, error) {
	err := i.repo.AddCollectionItem(ctx, request.GetBiz(), request.GetBizId(), request.GetCid(), request.GetUid())
	return &intrRepov1.AddCollectionItemResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) Get(ctx context.Context, request *intrRepov1.GetRequest) (*intrRepov1.GetResponse, error) {
	intr, err := i.repo.Get(ctx, request.GetBiz(), request.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrRepov1.GetResponse{
		Intr: i.toDTO(intr),
	}, err
}

func (i *InteractiveGRPCRepositoryServer) Liked(ctx context.Context, request *intrRepov1.LikedRequest) (*intrRepov1.LikedResponse, error) {
	_, err := i.repo.Liked(ctx, request.GetBiz(), request.GetId(), request.GetUid())
	return &intrRepov1.LikedResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) Collected(ctx context.Context, request *intrRepov1.CollectedRequest) (*intrRepov1.CollectedResponse, error) {
	_, err := i.repo.Collected(ctx, request.GetBiz(), request.GetId(), request.GetUid())
	return &intrRepov1.CollectedResponse{}, err
}

func (i *InteractiveGRPCRepositoryServer) GetByIds(ctx context.Context, request *intrRepov1.GetByIdsRequest) (*intrRepov1.GetByIdsResponse, error) {
	intrs, err := i.repo.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	res := make([]*intrRepov1.Interactive, len(intrs))
	for k, v := range intrs {
		res[k] = i.toDTO(v)
	}
	return &intrRepov1.GetByIdsResponse{
		Intrs: res,
	}, err
}

func NewInteractiveGRPCRepositoryServer(repo repository.InteractiveRepository) *InteractiveGRPCRepositoryServer{
	return &InteractiveGRPCRepositoryServer{
		repo: repo,
	}
}

func (i *InteractiveGRPCRepositoryServer) Register(server *grpc.Server) {
	intrRepov1.RegisterInteractiveRepositoryServer(server, i)
}

func (i *InteractiveGRPCRepositoryServer) toDTO(intr domain.Interactive) *intrRepov1.Interactive {
	return &intrRepov1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt,
	}
}

