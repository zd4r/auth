package app

import (
	"context"
	"net"

	"github.com/zd4r/auth/internal/config"
	"github.com/zd4r/auth/pkg/closer"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider

	grpcServer *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	var a App

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	err := a.runGRPCServer()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		config.Init,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = NewServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Host())
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}
