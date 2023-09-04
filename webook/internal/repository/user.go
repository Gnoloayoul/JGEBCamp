package repository

import (
	"context"
	"database/sql"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, dao.User{
		Id:       u.Id,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	})
}

func (r *UserRepository) Profile(ctx context.Context, u domain.User) (domain.User, error) {
	// 测试打印 取u之前
	//fmt.Printf("\nform repe before: %#v", u)
	resUser, err := r.dao.Profile(ctx, dao.User{
		Id: u.Id,
	})
	// 测试打印 取u之后
	//fmt.Printf("\nform repe after: %#v", resUser)
	return domain.User{
		Id:       resUser.Id,
		Password: resUser.Password,
		NickName: resUser.NickName,
		Birthday: resUser.Birthday,
		Info:     resUser.Info,
	}, err
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 第一次取，在缓存（redis）中
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	// 第二次取，在数据库（mysql）中
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志，做监控的好地方？
		}
	}()

	return u, err
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User {
		Id: u.Id,
		Email: u.Email.String,
		Password: u.Password,
		Phone: u.Phone.String,
		Ctime: time.UnixMilli(u.Ctime),
	}
}