package repository

import (
	"context"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, dao.User{
		Id: u.Id,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info: u.Info,
	})
}

func (r *UserRepository) Profile(ctx context.Context, u domain.User) (domain.User, error) {
	//return r.dao.Profile(ctx, dao.User{
	//	Id: u.Id,
	//})
	fmt.Printf("\nform repe before: %#v", u)
	resUser, err := r.dao.Profile(ctx, dao.User{
		Id: u.Id,
	})
	fmt.Printf("\nform repe after: %#v", resUser)
	return resUser ,err
}

func (r *UserRepository) FindById(int642 int64) {

}
