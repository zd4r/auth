package user

import (
	"context"
	"errors"
)

func (s *service) Delete(ctx context.Context, username string) error {
	rowsAffected, err := s.userRepository.Delete(ctx, username)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
