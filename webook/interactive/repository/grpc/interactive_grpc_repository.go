package grpc

import (
	intrRepov1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
)

type InteractiveGRPCRepositoryServer struct {
	intrRepov1.UnimplementedInteractiveRepositoryServer
	repo repository.InteractiveRepository
}



func NewInteractiveGRPCRepositoryServer(repo repository.InteractiveRepository) *InteractiveGRPCRepositoryServer{
	return &InteractiveGRPCRepositoryServer{
		repo: repo,
	}
}


