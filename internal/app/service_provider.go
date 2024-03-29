package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	userV1 "github.com/zd4r/auth/internal/api/user_v1"
	"github.com/zd4r/auth/internal/client/pg"
	"github.com/zd4r/auth/internal/config"
	userRepository "github.com/zd4r/auth/internal/repository/user"
	userService "github.com/zd4r/auth/internal/service/user"
	"github.com/zd4r/auth/pkg/closer"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	pgClient       pg.Client
	userRepository userRepository.Repository
	userService    userService.Service

	userImpl *userV1.Implementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetHTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) GetSwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) GetPgClient(ctx context.Context) pg.Client {
	if s.pgClient == nil {
		pgCfg, err := pgxpool.ParseConfig(s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to parse config from dsn")
		}

		cl, err := pg.NewClient(ctx, pgCfg)
		if err != nil {
			log.Fatalf("failed to get pg client: %s", err.Error())
		}

		err = cl.PG().Ping(ctx)
		if err != nil {
			log.Fatalf("failed ping db: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.pgClient = cl
	}

	return s.pgClient
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.GetPgClient(ctx))
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
