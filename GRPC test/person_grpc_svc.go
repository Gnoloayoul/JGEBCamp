package GRPC_test

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/GRPC test/api/proto/gen"
	"google.golang.org/grpc"
)

type PersonGRPCServer struct {
	gp PersonAction
	personv1.UnimplementedPersonActionServer
}

func NewPersonGRPCServer(gp PersonAction) *PersonGRPCServer {
	return &PersonGRPCServer{gp: gp}
}

func (p *PersonGRPCServer) Register(server *grpc.Server) {
	personv1.RegisterPersonActionServer(server, p)
}

func (p *PersonGRPCServer) SayHello(ctx context.Context, request *personv1.SayHelloRequest) (*personv1.SayHelloResponse, error) {
	p.gp.SayHello(request.GetAnybody())
	return &personv1.SayHelloResponse{}, nil
}

func (p *PersonGRPCServer) SayGoodBye(ctx context.Context, request *personv1.SayGoodByeRequest) (*personv1.SayGoodByeResponse, error) {
	p.gp.SayGoodBye(request.GetAnybody())
	return &personv1.SayGoodByeResponse{}, nil
}
