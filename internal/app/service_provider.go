package app

import (
	"context"
	"log"

	userV1 "github.com/zd4r/auth/internal/api/user_v1"
	"github.com/zd4r/auth/internal/client/pg"
	userRepository "github.com/zd4r/auth/internal/repository/user"
	userService "github.com/zd4r/auth/internal/service/user"
	"github.com/zd4r/auth/pkg/closer"
)

type serviceProvider struct {
	//pgxPool        *pgxpool.Pool
	pgClient       pg.Client
	userRepository userRepository.Repository
	userService    userService.Service

	userImpl *userV1.Implementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPgClient(ctx context.Context) pg.Client {
	if s.pgClient == nil {
		cl, err := pg.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to get pg client: %s", err.Error())
		}

		err = cl.PG().Ping(ctx)
		if err != nil {
			log.Fatalf("failed ping db: %s", err.Error())
		}
		closer.Add(cl.PG().Close)

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
