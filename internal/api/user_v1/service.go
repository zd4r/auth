package user_v1

import (
	"github.com/zd4r/auth/internal/service/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server

	userService user.Service
}

func NewImplementation(userService user.Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
