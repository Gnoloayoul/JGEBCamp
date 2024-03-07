package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/grpcx"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
