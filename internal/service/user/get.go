package user

import (
	"context"

	model "github.com/zd4r/auth/internal/model/user"
)

func (s *service) Get(ctx context.Context, username string) (*model.User, error) {
	u, err := s.userRepository.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	return u, nil
}
