package config

import (
	"errors"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
)

var _ GRPCConfig = (*grpcConfig)(nil)

type GRPCConfig interface {
	Host() string
}

type grpcConfig struct {
	host string
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	return &grpcConfig{
		host: host,
	}, nil
}

func (cfg *grpcConfig) Host() string {
	return cfg.host
}
