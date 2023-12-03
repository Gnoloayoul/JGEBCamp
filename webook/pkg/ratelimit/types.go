package ratelimit

import "context"

//go:generate mockgen -source=./types.go -package=ratelimitmocks -destination=mocks/types_mock.go Limiter
type Limiter interface {
	Limit(ctx context.Context, key string) (bool, error)
}
