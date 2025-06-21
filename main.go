package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	accountv1 "github.com/0utl1er-tech/go-backend/gen/pb/account/service/v1grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/0utl1er-tech/go-backend/internal/application/mail"
	"github.com/0utl1er-tech/go-backend/internal/config"
	"github.com/0utl1er-tech/go-backend/internal/gapi/method"
	"github.com/0utl1er-tech/go-backend/internal/gapi/service"
	"github.com/0utl1er-tech/go-backend/internal/infra/db/repository"
	"github.com/0utl1er-tech/go-backend/internal/infra/worker"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hibiken/asynq"
	"github.com/rakyll/statik/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if config.Environment == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db")
	}

	store := repository.NewDAO(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr:     net.JoinHostPort(config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Config,
	redisOpt asynq.RedisClientOpt,
	store repository.DAO,
) {
	mailer := *mail.NewSendConfirmationMail(config.EmailSender.Name, config.EmailSender.Address, config.EmailSender.Password)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("Start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("Task processor is stopped")

		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Config,
	store repository.DAO,
	taskDistributor worker.TaskDistributor,
) {
	server, err := service.NewServer(config, store, taskDistributor)
	method := method.NewMethod(server)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	gprcLogger := grpc.UnaryInterceptor(service.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	accountv1.RegisterAccountServiceServer(grpcServer, method)
	reflection.Register(grpcServer)

	GRPCServerAddress := net.JoinHostPort(config.Server.Host, config.Server.GRPCPort)
	listener, err := net.Listen("tcp", GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener")
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

	// graceful shutdown
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
	config config.Config,
	store repository.DAO,
	taskDistributor worker.TaskDistributor,
) {
	server, err := service.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	err = accountv1.RegisterAccountServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	})
	handler := c.Handler(service.HttpLogger(mux))

	httpServer := &http.Server{
		Handler: handler,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
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
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}

		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}
