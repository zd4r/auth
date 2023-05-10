package user

import (
	"context"

	model "github.com/zd4r/auth/internal/model/user"
)

func (s *service) Get(ctx context.Context, user *model.User) (*model.User, error) {
	u, err := s.userRepository.Get(ctx, user)
	if err != nil {
		return nil, err
	}

	return u, nil
}
