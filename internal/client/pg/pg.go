package pg

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ PG = (*pg)(nil)

type Query struct {
	Name     string
	QueryRaw string
}

type QueryExecer interface {
	Exec(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type NamedScanner interface {
	ScanOne(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAll(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Closer interface {
	Close()
}

type PG interface {
	QueryExecer
	NamedScanner
	Pinger
	Closer
}

type pg struct {
	pgxPool *pgxpool.Pool
}

func (p *pg) ScanOne(ctx context.Context, dest interface{}, q Query, args ...interface{}) error {
	row, err := p.Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAll(ctx context.Context, dest interface{}, q Query, args ...interface{}) error {
	rows, err := p.Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
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

func (p *pg) Ping(ctx context.Context) error {
	return p.pgxPool.Ping(ctx)
}

func (p *pg) Close() {
	p.pgxPool.Close()
}
