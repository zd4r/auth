package user

import (
	"database/sql"

	model "github.com/zd4r/auth/internal/model/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToUser(user *desc.User) *model.User {
	u := &model.User{
		Role: user.GetRole().String(),
	}

	if user.GetUsername() != nil {
		u.Username = sql.NullString{
			String: user.GetUsername().GetValue(),
			Valid:  true,
		}
	}

	if user.GetPassword() != nil {
		u.Password = sql.NullString{
			String: user.GetPassword().GetValue(),
			Valid:  true,
		}
	}

	if user.GetPasswordConfirm() != nil {
		u.ConfirmPassword = sql.NullString{
			String: user.GetPasswordConfirm().GetValue(),
			Valid:  true,
		}
	}

	if user.GetEmail() != nil {
		u.Email = sql.NullString{
			String: user.GetEmail().GetValue(),
			Valid:  true,
		}
	}

	return u
}

func ToUsername(username string) *model.User {
	return &model.User{
		Username: sql.NullString{
			String: username,
			Valid:  true,
		},
	}
}

func ToUserDesc(user *model.User) *desc.UserInfo {
	u := desc.UserInfo{
		User: &desc.User{},
	}

	if user.Username.Valid {
		u.User.Username = wrapperspb.String(user.Username.String)
	}

	if user.Password.Valid {
		u.User.Password = wrapperspb.String(user.Password.String)
	}

	if user.Email.Valid {
		u.User.Email = wrapperspb.String(user.Email.String)
	}

	// Верно ли??
	u.User.Role = desc.RoleInfo(desc.RoleInfo_value[user.Role])

	u.CreatedAt = timestamppb.New(user.CreatedAt)

	u.UpdatedAt = timestamppb.New(user.UpdatedAt)

	return &u
}
