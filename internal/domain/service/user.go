package service

import (
	"context"
	"webTemplate/internal/domain/common/errorz"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/entity"
)

type userStorage interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type UserService struct {
	storage userStorage
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) Create(ctx context.Context, registerReq dto.UserRegister) (*entity.User, error) {
	if _, err := s.storage.GetByEmail(ctx, registerReq.Email); err == nil {
		return nil, errorz.EmailAlreadyExists
	}

	user := entity.User{
		Email: registerReq.Email,
	}
	user.SetPassword(registerReq.Password)
	return s.storage.Create(ctx, user)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.storage.GetByEmail(ctx, email)
}
