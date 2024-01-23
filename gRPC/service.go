package gRPC

import "context"

type Server struct {
	UnimplementedUserSvcServer
}

var _ UserSvcServer = &Server{}

func (s *Server) GetById(ctx context.Context, request *GetByIdRep) (*GetByIdResp, error) {
	return &GetByIdResp{
		User: &User{
			Id: 123,
			Name: "abcd",
		},
	}, nil
}