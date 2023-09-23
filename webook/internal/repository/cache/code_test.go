package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCodeRedisCache(tc.mock(ctrl))
			c.Set(tc.ctx, )
		})
	}
}
