package user

import (
	"context"
	"errors"

	model "github.com/zd4r/auth/internal/model/user"
)

func (s *service) Create(ctx context.Context, user *model.User) error {
	if user.Password.String != user.ConfirmPassword {
		return errors.New("passwords don't match")
	}

	err := s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
