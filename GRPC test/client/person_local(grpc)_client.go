package client

import (
	"context"
	GRPC_test "github.com/Gnoloayoul/JGEBCamp/GRPC test"
	"github.com/Gnoloayoul/JGEBCamp/GRPC test/api/proto/gen"
	"google.golang.org/grpc"
)

type PersonLocalAdapter struct {
	gp GRPC_test.PersonAction
}

func (p *PersonLocalAdapter) SayHello(ctx context.Context, in *personv1.SayHelloRequest, opts ...grpc.CallOption) (*personv1.SayHelloResponse, error) {
	p.gp.SayHello(in.GetAnybody())
	return &personv1.SayHelloResponse{}, nil
}

func (p *PersonLocalAdapter) SayGoodBye(ctx context.Context, in *personv1.SayGoodByeRequest, opts ...grpc.CallOption) (*personv1.SayGoodByeResponse, error) {
	p.gp.SayGoodBye(in.GetAnybody())
	return &personv1.SayGoodByeResponse{}, nil
}

func NewPersonLocalAdapter(gp GRPC_test.PersonAction) *PersonLocalAdapter {
	return &PersonLocalAdapter{
		gp: gp,
	}
}
