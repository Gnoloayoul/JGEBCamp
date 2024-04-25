package client

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/GRPC test/api/proto/gen"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"math/rand"
)

type PersonGreyScaleClient struct {
	// 本地
	local *PersonLocalAdapter
	// grpc
	remote personv1.PersonActionClient
	// 阈值
	threshold *atomicx.Value[int32]
}

func (p *PersonGreyScaleClient) client() personv1.PersonActionClient {
	threshold := p.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return p.remote
	}
	return p.local
}

func (p *PersonGreyScaleClient) SayHello(ctx context.Context, in *personv1.SayHelloRequest, opts ...grpc.CallOption) (*personv1.SayHelloResponse, error) {
	return p.client().SayHello(ctx, in, opts...)
}

func (p *PersonGreyScaleClient) SayGoodBye(ctx context.Context, in *personv1.SayGoodByeRequest, opts ...grpc.CallOption) (*personv1.SayGoodByeResponse, error) {
	return p.client().SayGoodBye(ctx, in, opts...)
}

func NewPersonGreyScaleClient(local *PersonLocalAdapter, remote personv1.PersonActionClient) *PersonGreyScaleClient {
	return &PersonGreyScaleClient{
		local:     local,
		remote:    local,
		threshold: atomicx.NewValue[int32](),
	}
}
