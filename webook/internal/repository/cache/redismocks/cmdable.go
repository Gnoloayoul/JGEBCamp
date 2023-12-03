package redismocks

//go:generate mockgen -package=redismocks -destination=./cmdable_mock.go github.com/redis/go-redis/v9 Cmdable
