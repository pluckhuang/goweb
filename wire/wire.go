//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/wire/repository"
	"github.com/pluckhuang/goweb/wire/repository/dao"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, InitDB, dao.NewUserDAO)
	return &repository.UserRepository{}
}
