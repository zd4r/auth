package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Client = (*client)(nil)

type Client interface {
	Close() error
	PG() PG
}

type client struct {
	pg PG
}

func NewClient(ctx context.Context, pgCfg *pgxpool.Config) (*client, error) {
	dbc, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		log.Fatalf("failed to get db connections: %s", err.Error())
	}

	return &client{
		pg: &pg{
			pgxPool: dbc,
		},
	}, nil
}

func (c *client) PG() PG {
	return c.pg
}

// Close - Зачем если можно Client.PG().Close() ? (также как и с Client.PG().Ping())
func (c *client) Close() error {
	if c.pg != nil {
		return c.pg.Close()
	}

	return nil
}
