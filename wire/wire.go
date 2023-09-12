//go:build wireinject

package wire

import (
	"github.com/Gnoloayoul/JGEBCamp/wire/repository"
	"github.com/Gnoloayoul/JGEBCamp/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	wire.Build(
		repository.NewUserRepository,
		dao.NewUserDAO,
		InitDB)
	return new(repository.UserRepository)
}
