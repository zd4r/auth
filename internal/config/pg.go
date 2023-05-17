package config

import (
	"errors"
	"os"
)

const (
	pgDSNEnvName = "PG_DSN"
)

var _ PGConfig = (*pgConfig)(nil)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(pgDSNEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
