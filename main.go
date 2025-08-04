package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"syscall"

	customerv1 "github.com/0utl1er-tech/prism-backend/gen/pb/customer/v1"
	db "github.com/0utl1er-tech/prism-backend/gen/sqlc"
	"github.com/0utl1er-tech/prism-backend/internal/service"
	"github.com/0utl1er-tech/prism-backend/internal/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	connPool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create connection pool")
	}

	queries := db.New(connPool)
	customerService := service.NewCustomerService(queries)

	waitGroup, ctx := errgroup.WithContext(context.Background())
	// TODO: 引数がデカすぎるのでリファクタリングする
	runGrpcServer(ctx, waitGroup, customerService, &cfg)
	runGatewayServer(ctx, waitGroup, customerService, &cfg)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to wait")
	}

}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	customer *service.CustomerService,
	cfg *util.Config,
) {
	grpcServer := grpc.NewServer()

	customerv1.RegisterCustomerServiceServer(grpcServer, customer)

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	waitGroup.Go(func() error {
		log.Info().Msgf("Start gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")
		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	customer *service.CustomerService,
	cfg *util.Config,
) {
	// grpc-ecosystemのmiddlewareを使用したServeMuxオプション
	serveMuxOptions := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}

	serverMuxOptions := append(serveMuxOptions)
	grpcMux := runtime.NewServeMux(serverMuxOptions...)

	err := customerv1.RegisterCustomerServiceHandlerServer(ctx, grpcMux, customer)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register customer service handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	httpServer := &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: mux,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("Start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown HTTP gateway server")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("HTTP gateway server shutdown error")
			return err
		}
		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}
