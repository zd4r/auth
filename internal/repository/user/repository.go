package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zd4r/auth/internal/client/pg"
	model "github.com/zd4r/auth/internal/model/user"
)

var _ Repository = (*repository)(nil)

const tableName = `"user"`

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, username string, user *model.User) (int64, error)
	Delete(ctx context.Context, username string) (int64, error)
}

type repository struct {
	client pg.Client
}

func NewRepository(client pg.Client) *repository {
	return &repository{
		client: client,
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

	q := pg.Query{
		Name:     "user.Create",
		QueryRaw: query,
	}

	_, err = r.client.PG().Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, username string) (*model.User, error) {
	builder := sq.Select("username", "email", "password", "role", "created_at", "updated_at").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"username": username,
		}).Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "user.Get",
		QueryRaw: query,
	}

	var u model.User
	err = r.client.PG().ScanOne(ctx, &u, q, args...)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) Update(ctx context.Context, username string, user *model.User) (int64, error) {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"username": username,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := pg.Query{
		Name:     "user.Update",
		QueryRaw: query,
	}

	ct, err := r.client.PG().Exec(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}

func (r *repository) Delete(ctx context.Context, username string) (int64, error) {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"username": username,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := pg.Query{
		Name:     "user.Delete",
		QueryRaw: query,
	}

	ct, err := r.client.PG().Exec(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}
