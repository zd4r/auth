package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	model "github.com/zd4r/auth/internal/model/user"
)

var _ Repository = (*repository)(nil)

const tableName = `"user"`

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns("username", "email", "password", "role").
		Values(user.Username, user.Email, user.Password, user.Role)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, user *model.User) (*model.User, error) {
	builder := sq.Select("username", "email", "password", "role", "created_at", "updated_at").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"username": user.Username,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	// rows.Next() automatically closes rows when all rows are read (тогда зачем rows.Close()?)
	defer rows.Close()

	// Нет цикла, поскольку предполагается что username уникальный.
	rows.Next()
	var u model.User
	err = rows.Scan(
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) Update(ctx context.Context, user *model.User) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"username": user.Username,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, user *model.User) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"username": user.Username,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
