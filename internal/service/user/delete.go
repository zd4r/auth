package user

import (
	"context"

	model "github.com/zd4r/auth/internal/model/user"
)

func (s *service) Delete(ctx context.Context, user *model.User) error {
	err := s.userRepository.Delete(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
