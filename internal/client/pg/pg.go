package pg

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Query struct {
	Name     string
	QueryRaw string
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Closer interface {
	Close() error
}

type QueryExecer interface {
	Exec(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type PG interface {
	Pinger
	Closer
	QueryExecer
}

type pg struct {
	pgxPool *pgxpool.Pool
}

func (p *pg) Ping(ctx context.Context) error {
	return p.pgxPool.Ping(ctx)
}

func (p *pg) Close() error {
	p.pgxPool.Close()
	return nil
}

func (p *pg) Exec(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error) {
	return p.pgxPool.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) Query(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error) {
	return p.pgxPool.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRow(ctx context.Context, q Query, args ...interface{}) pgx.Row {
	return p.pgxPool.QueryRow(ctx, q.QueryRaw, args...)
}
