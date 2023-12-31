package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁了")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知道发生了什么，反正就是和 Code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

//go:generate mockgen -source=./code.go -package=cachemocks -destination=mocks/code_mock.go CodeCache
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeRedisCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
		//case -2:
	//	return
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

//func (c *CodeCache) Verify(ctx context.Context, biz, phone, code string) error {
//
//}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// LocalCodeCache 假如说你要切换这个，你是不是得把 lua 脚本的逻辑，在这里再写一遍？
type LocalCodeCache struct {
	client redis.Cmdable
}

// 本地缓存版
type CodeLocalCache struct {
	data sync.Map
}

func NewCodeCache(data sync.Map) *CodeLocalCache {
	return &CodeLocalCache{
		data: data,
	}
}

func (c *CodeLocalCache) Set(biz, phone, code string) error {
	_, ok := c.data.Load(biz + phone)
	if !ok {
		c.data.Store(biz+phone, code)
	}
	return nil
}

func (c *CodeLocalCache) Verify(biz, phone, inputCode string) (bool, error) {
	curCode, ok := c.data.Load(biz + phone)
	if !ok {
		// 该电话号码还没配备验证码
		return false, errors.New("系统错误")
	}
	if curCode != inputCode {
		// 输入的验证错误
		return false, nil
	}
	return true, nil
}

func (c *CodeLocalCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
