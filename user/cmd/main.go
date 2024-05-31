package main

import (
	"context"
	repositoryFriendInvte "github.com/Nixonxp/discord/user/internal/app/repository/friend_invites_storage"
	repositoryUserFriends "github.com/Nixonxp/discord/user/internal/app/repository/user_friends_storage"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/user/internal/app/server"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/user/internal/middleware/errors"
	middleware_tracing "github.com/Nixonxp/discord/user/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	"github.com/Nixonxp/discord/user/pkg/application"
	logCfg "github.com/Nixonxp/discord/user/pkg/logger"
	logger "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/Nixonxp/discord/user/pkg/postgres"
	"github.com/Nixonxp/discord/user/pkg/rate_limiter"
	jaeger_tracing "github.com/Nixonxp/discord/user/pkg/tracing"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
	"time"
)

const DSN = "user=admin password=password123 host=postgres port=5432 dbname=discord sslmode=require pool_max_conns=10"

//const DSN = "user=admin password=password123 host=localhost port=5432 dbname=discord sslmode=require pool_max_conns=10" // todo delete

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer resourcesShutdownCtxCancel()

	config := application.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
	}

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	if err := jaeger_tracing.Init("user service"); err != nil {
		log.Fatal(ctx, err)
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	// repository
	pool, err := postgres.NewConnectionPool(resourcesShutdownCtx, DSN,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		log.Fatal(err)
	}

	txManager := transaction_manager.New(pool)

	userRepo := repository.NewUserPostgresqlRepository(txManager, log)
	friendInvitesRepo := repositoryFriendInvte.NewFriendInvitesPostgresqlRepository(txManager, log)
	friendsRepo := repositoryUserFriends.NewUserFriendsPostgresqlRepository(txManager, log)

	userUsecase := usecases.NewUserUsecase(usecases.Deps{
		UserRepo:           userRepo,
		FriendInvitesRepo:  friendInvitesRepo,
		UserFriendsRepo:    friendsRepo,
		TransactionManager: txManager,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(log),
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
		},
	}

	srv, err := server.NewUserServer(resourcesShutdownCtx, server.Deps{
		UserUsecase: userUsecase,
		Log:         log,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterUserServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")
	resourcesShutdownCtxCancel()
	log.Print("wait shutdown resources")
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
