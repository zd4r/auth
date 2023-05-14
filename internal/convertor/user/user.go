package user

import (
	"database/sql"

	model "github.com/zd4r/auth/internal/model/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUser(user *desc.User, passwordConfirm string) *model.User {
	u := model.User{}

	u.Username = sql.NullString{
		String: user.GetUsername(),
		Valid:  true,
	}

	u.Email = sql.NullString{
		String: user.GetEmail(),
		Valid:  true,
	}

	u.Password = sql.NullString{
		String: user.GetPassword(),
		Valid:  true,
	}

	u.ConfirmPassword = sql.NullString{
		String: passwordConfirm,
		Valid:  true,
	}

	u.Role = user.GetRole().String()

	return &u
}

func ToUserNullable(user *desc.UserNullable) *model.User {
	u := model.User{}

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

	if user.GetEmail() != nil {
		u.Email = sql.NullString{
			String: user.GetEmail().GetValue(),
			Valid:  true,
		}
	}

	u.Role = user.GetRole().String()

	return &u
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
		u.User.Username = user.Username.String
	}

	if user.Password.Valid {
		u.User.Password = user.Password.String
	}

	if user.Email.Valid {
		u.User.Email = user.Email.String
	}

	u.User.Role = desc.RoleInfo(desc.RoleInfo_value[user.Role])

	u.CreatedAt = timestamppb.New(user.CreatedAt)

	u.UpdatedAt = timestamppb.New(user.UpdatedAt)

	return &u
}
