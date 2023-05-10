package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	userV1 "github.com/zd4r/auth/internal/api/user_v1"
	userRepository "github.com/zd4r/auth/internal/repository/user"
	userService "github.com/zd4r/auth/internal/service/user"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50051"

func main() {
	ctx := context.Background()

	list, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to get listener: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pgCfg, err := pgxpool.ParseConfig("host=localhost port=54321 dbname=user user=user-user password=user-password sslmode=disable")
	if err != nil {
		log.Fatalf("failed to parse config from dsn")
	}
	dbc, err := pgxpool.ConnectConfig(ctx, pgCfg)
	if err != nil {
		log.Fatalf("failed to get db connections: %s", err.Error())
	}
	// Не выполниться из-за log.Fatalf
	defer dbc.Close()

	err = dbc.Ping(ctx)
	if err != nil {
		log.Fatalf("failed ping db: %s", err.Error())
	}

	userRepo := userRepository.NewRepository(dbc)
	userSrv := userService.NewService(userRepo)
	desc.RegisterUserV1Server(s, userV1.NewImplementation(userSrv))

	err = s.Serve(list)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
