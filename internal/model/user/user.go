package user

import (
	"database/sql"
	"time"
)

type User struct {
	Username        sql.NullString
	Email           sql.NullString
	Password        sql.NullString
	ConfirmPassword sql.NullString
	Role            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
