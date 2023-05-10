package user_v1

import (
	"context"

	convertor "github.com/zd4r/auth/internal/convertor/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	u, err := i.userService.Get(ctx, convertor.ToUsername(req.GetUsername()))
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		UserInfo: convertor.ToUserDesc(u),
	}, nil
}
