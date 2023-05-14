package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	userV1 "github.com/zd4r/auth/internal/api/user_v1"
	userRepository "github.com/zd4r/auth/internal/repository/user"
	userService "github.com/zd4r/auth/internal/service/user"
	"github.com/zd4r/auth/pkg/closer"
)

type serviceProvider struct {
	pgxPool        *pgxpool.Pool
	userRepository userRepository.Repository
	userService    userService.Service

	userImpl *userV1.Implementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPgxPool(ctx context.Context) *pgxpool.Pool {
	if s.pgxPool == nil {
		pgCfg, err := pgxpool.ParseConfig("host=localhost port=54321 dbname=user user=user-user password=user-password sslmode=disable")
		if err != nil {
			log.Fatalf("failed to parse config from dsn")
		}

		dbc, err := pgxpool.ConnectConfig(ctx, pgCfg)
		if err != nil {
			log.Fatalf("failed to get db connections: %s", err.Error())
		}
		closer.Add(func() error {
			dbc.Close()
			return nil
		})

		err = dbc.Ping(ctx)
		if err != nil {
			log.Fatalf("failed ping db: %s", err.Error())
		}

		s.pgxPool = dbc
	}

	return s.pgxPool
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.GetPgxPool(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) GetUserService(ctx context.Context) userService.Service {
	if s.userService == nil {
		s.userService = userService.NewService(s.GetUserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) GetUserImpl(ctx context.Context) *userV1.Implementation {
	if s.userImpl == nil {
		s.userImpl = userV1.NewImplementation(s.GetUserService(ctx))
	}

	return s.userImpl
}
