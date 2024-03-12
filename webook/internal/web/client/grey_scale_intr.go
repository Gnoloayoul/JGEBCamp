package client

import (
	"context"
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"math/rand"
)

// GreyScaleInteractiveServiceClient
// 运用阈值+随机数控制流量的 Client
type GreyScaleInteractiveServiceClient struct {
	remote intrv1.InteractiveServiceClient
	local *InteractiveServiceAdapter
	threshold *atomicx.Value[int32]
}

func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return g.client().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return g.client().GetByIds(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) OnChange (ch <-chan int32) {
	go func() {
		for newch := range ch {
			g.threshold.Store(newch)
		}
	}()
}

func (g *GreyScaleInteractiveServiceClient) OnChangeV1() chan<- int32 {
	ch := make(chan int32, 100)
	go func() {
		for newCh := range ch {
			g.threshold.Store(newCh)
		}
	}()
	return ch
}

func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(newThreshold int32) {
	g.threshold.Store(newThreshold)
}

// client
// 在内部生成随机数，与阈值比对，从而决定流量去向
func (g *GreyScaleInteractiveServiceClient) client() intrv1.InteractiveServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.remote
	}
	return g.local
}

func NewGreyScaleInteractiveServiceClient(remote intrv1.InteractiveServiceClient, local *InteractiveServiceAdapter) *GreyScaleInteractiveServiceClient{
	return &GreyScaleInteractiveServiceClient{
		remote: remote,
		local: local,
		threshold: atomicx.NewValue[int32()](),
	}
}

