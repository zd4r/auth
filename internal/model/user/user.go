package user

import (
	"database/sql"
	"time"
)

type User struct {
	Username        sql.NullString
	Email           sql.NullString
	Password        sql.NullString
	ConfirmPassword string
	Role            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
