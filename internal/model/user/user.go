package user

import (
	"database/sql"
	"time"
)

type User struct {
	Username        sql.NullString `db:"username"`
	Email           sql.NullString `db:"email"`
	Password        sql.NullString `db:"password"`
	ConfirmPassword string
	Role            string    `db:"role"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
