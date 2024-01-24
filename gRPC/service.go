package gRPC

import (
	"context"
	"errors"
)

type Server struct {
	UnimplementedUserSvcServer
}

var _ UserSvcServer = &Server{}

//func (s *Server) GetById(ctx context.Context, request *GetByIdRep) (*GetByIdResp, error) {
//	return &GetByIdResp{
//		User: &User{
//			Id: 123,
//			Name: "abcd",
//		},
//	}, nil
//}

func (s *Server) GetById(ctx context.Context, request *GetByIdRep) (*GetByIdResp, error) {
	return check(request)
}

func check(request *GetByIdRep) (*GetByIdResp, error) {
	//switch request.Id {
	//case 123:
	//	return &GetByIdResp{
	//		User: &User{
	//			Id: 123,
	//			Name: "abcd",
	//		},
	//	}, nil
	//case 456:
	//	return &GetByIdResp{
	//		User: &User{
	//			Id: 456,
	//			Name: "AAAA",
	//		},
	//	}, nil
	//default:
	//	return &GetByIdResp{User: &User{}}, errors.New("bad req!")
	//}
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
		return &GetByIdResp{User: &User{}}, errors.New("bad req!")
	} else {
		return list[request.Id], nil
	}
}
