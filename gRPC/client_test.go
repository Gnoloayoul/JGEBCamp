package gRPC

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cc, err := grpc.Dial(":8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := NewUserSvcClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdRep{
		Id: 333,
	})
	assert.NoError(t, err)
	t.Log(resp.User)
}

