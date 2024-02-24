package grpc

import (
	"context"
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}


