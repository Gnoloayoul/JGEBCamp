package client

import (
	"context"
	intrRepov1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyscaleInteractiveRepositoryClient struct {
	remote    intrRepov1.InteractiveRepositoryClient
	local     *InteractiveRepositoryAdapter
	threshold *atomicx.Value[int32]
}

func (g *GreyscaleInteractiveRepositoryClient) client() intrRepov1.InteractiveRepositoryClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.remote
	}
	return g.local
}

func (g *GreyscaleInteractiveRepositoryClient) IncrReadCnt(ctx context.Context, in *intrRepov1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrRepov1.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) BatchIncrReadCnt(ctx context.Context, in *intrRepov1.BatchIncrReadCntRequest, opts ...grpc.CallOption) (*intrRepov1.BatchIncrReadCntResponse, error) {
	return g.client().BatchIncrReadCnt(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) IncrLike(ctx context.Context, in *intrRepov1.IncrLikeRequest, opts ...grpc.CallOption) (*intrRepov1.IncrLikeResponse, error) {
	return g.client().IncrLike(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) DecrLike(ctx context.Context, in *intrRepov1.DecrLikeRequest, opts ...grpc.CallOption) (*intrRepov1.DecrLikeResponse, error) {
	return g.client().DecrLike(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) AddCollectionItem(ctx context.Context, in *intrRepov1.AddCollectionItemRequest, opts ...grpc.CallOption) (*intrRepov1.AddCollectionItemResponse, error) {
	return g.client().AddCollectionItem(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) Get(ctx context.Context, in *intrRepov1.GetRequest, opts ...grpc.CallOption) (*intrRepov1.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) Liked(ctx context.Context, in *intrRepov1.LikedRequest, opts ...grpc.CallOption) (*intrRepov1.LikedResponse, error) {
	return g.client().Liked(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) Collected(ctx context.Context, in *intrRepov1.CollectedRequest, opts ...grpc.CallOption) (*intrRepov1.CollectedResponse, error) {
	return g.client().Collected(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) GetByIds(ctx context.Context, in *intrRepov1.GetByIdsRequest, opts ...grpc.CallOption) (*intrRepov1.GetByIdsResponse, error) {
	return g.client().GetByIds(ctx, in, opts...)
}

func (g *GreyscaleInteractiveRepositoryClient) UpdateThreshold(newThreshold int32) {
	g.threshold.Store(newThreshold)
}

func NewGreyscaleInteractiveRepositoryClient(remote intrRepov1.InteractiveRepositoryClient, local *InteractiveRepositoryAdapter) *GreyscaleInteractiveRepositoryClient {
	return &GreyscaleInteractiveRepositoryClient{
		remote: remote, local: local,
		threshold: atomicx.NewValue[int32()](),
	}
}
