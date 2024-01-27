package gRPC

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedUserSvcServer
}

var _ UserSvcServer = &Server{}

func (s *Server) GetById(ctx context.Context, request *GetByIdRep) (*GetByIdResp, error) {
		list := map[int64]*GetByIdResp {
			123: &GetByIdResp{User: &User{
				Id: 123,
				Name: "abcd",
			}},
			456: &GetByIdResp{User: &User{
				Id: 456,
				Name: "AAAA",
			}},
		}

	_, ok := list[request.Id]
	if !ok {
		// 创建gRPC错误
		err := status.Errorf(codes.NotFound, "没有该用户")
		return nil, err
	}

	return list[request.Id], nil
}
