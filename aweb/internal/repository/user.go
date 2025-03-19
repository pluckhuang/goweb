package repository

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/internal/domain"
	"github.com/pluckhuang/goweb/aweb/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) UpdateInfoById(ctx context.Context, userId int64, updateFields map[string]interface{}) error {
	return repo.dao.UpdateById(ctx, userId, updateFields)
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:          u.Id,
		Email:       u.Email,
		Password:    u.Password,
		Birthday:    u.Birthday.UnixMilli(),
		Description: u.Description,
		Nickname:    u.Nickname,
	}
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
