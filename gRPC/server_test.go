package gRPC

import (
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer()
	defer func() {
		server.GracefulStop()
	}()

	userServer := &Server{}
	RegisterUserSvcServer(server, userServer)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	t.Log(err)
}