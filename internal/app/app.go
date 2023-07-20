package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zd4r/auth/internal/config"
	"github.com/zd4r/auth/internal/interceptor"
	"github.com/zd4r/auth/pkg/closer"
	desc "github.com/zd4r/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider

	grpcServer *grpc.Server
	httpServer *http.Server
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

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := a.runGRPCServer(); err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		config.Init,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := desc.RegisterUserV1HandlerFromEndpoint(ctx,
		mux,
		a.serviceProvider.GetGRPCConfig().Host(),
		opts,
	); err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.GetHTTPConfig().Host(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s\n", a.serviceProvider.GetGRPCConfig().Host())

	lis, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Host())
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s\n", a.serviceProvider.GetHTTPConfig().Host())

	return a.httpServer.ListenAndServe()
}
