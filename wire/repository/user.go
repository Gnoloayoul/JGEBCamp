package repository

import "github.com/Gnoloayoul/JGEBCamp/wire/repository/dao"

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}