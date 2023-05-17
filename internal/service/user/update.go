package user

import (
	"context"

	model "github.com/zd4r/auth/internal/model/user"
)

func (s *service) Update(ctx context.Context, username string, user *model.User) error {
	latestUser, err := s.userRepository.Get(ctx, username)
	if err != nil {
		return err
	}

	match(user, latestUser)

	err = s.userRepository.Update(ctx, username, latestUser)
	if err != nil {
		return err
	}

	return nil
}

func match(user, latestUser *model.User) {
	if user.Username.Valid {
		latestUser.Username = user.Username
	}

	if user.Email.Valid {
		latestUser.Email = user.Email
	}

	if user.Password.Valid {
		latestUser.Password = user.Password
	}

	latestUser.Role = user.Role
}
