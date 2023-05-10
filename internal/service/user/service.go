package user

import (
	"context"

	model "github.com/zd4r/auth/internal/model/user"
	"github.com/zd4r/auth/internal/repository/user"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
}

type service struct {
	userRepository user.Repository
}

func NewService(userRepository user.Repository) *service {
	return &service{
		userRepository: userRepository,
	}
}
