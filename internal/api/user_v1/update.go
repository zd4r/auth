package user_v1

import (
	"context"

	convertor "github.com/zd4r/auth/internal/convertor/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req.GetUsername(), convertor.ToUserNullable(req.GetUser()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
