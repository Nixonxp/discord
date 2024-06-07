package main

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/repository/server_storage"
	subscribe_storage "github.com/Nixonxp/discord/server/internal/app/repository/subscribe_storage"
	"github.com/Nixonxp/discord/server/internal/app/server"
	"github.com/Nixonxp/discord/server/internal/app/services/chat"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/server/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/server/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/server/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/server/pkg/api/v1"
	"github.com/Nixonxp/discord/server/pkg/application"
	logCfg "github.com/Nixonxp/discord/server/pkg/logger"
	logger "github.com/Nixonxp/discord/server/pkg/logger"
	mongoCollection "github.com/Nixonxp/discord/server/pkg/mongo"
	"github.com/Nixonxp/discord/server/pkg/rate_limiter"
	jaeger_tracing "github.com/Nixonxp/discord/server/pkg/tracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	ratelimitCustom "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer resourcesShutdownCtxCancel()

	config := application.Config{
		GRPCPort: ":8480", // todo :8080
		//HTTPPort:  ":8081",
		//DebugPort: ":8084",
	}

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	closer, err := jaeger_tracing.Init("server service")
	if err != nil {
		log.Fatal(ctx, err)
	}
	defer closer(resourcesShutdownCtx)

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	serverConfig := server.Config{
		ChatServiceUrl: "localhost:8580", // todo chat:8080
		UnaryClientInterceptors: []grpc.UnaryClientInterceptor{
			grpcretry.UnaryClientInterceptor(retryOpts...),
			ratelimitCustom.UnaryClientInterceptor(ratelimitCustom.NewLimiter(10001)),
			grpc_opentracing.OpenTracingClientInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
		},
	}

	chatServiceClient, err := chat.NewClient(serverConfig, log)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	collection, err := mongoCollection.NewCollection(resourcesShutdownCtx,
		"servers",
		&mongoCollection.Config{
			/*MongoHost:     os.Getenv("MONGO_HOST"),
			MongoDb:       os.Getenv("MONGO_DB"),
			MongoPort:     os.Getenv("MONGO_PORT"),
			MongoUser:     os.Getenv("MONGO_USER"),
			MongoPassword: os.Getenv("MONGO_PASSWORD"),*/

			MongoHost:     "localhost",
			MongoDb:       "discord",
			MongoPort:     "27117",
			MongoUser:     "discord",
			MongoPassword: "example",
			// todo local vars delete
		},
	)

	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	subscribeCollection, err := collection.NewCollection("server_subscribe")
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	serverMongoRepo := repository.NewMongoServerRepository(collection, log)
	subscribeMongoRepo := subscribe_storage.NewMongoSubscribeRepository(subscribeCollection, log)

	serverUsecase := usecases.NewServerUsecase(usecases.Deps{
		ServerRepo:    serverMongoRepo,
		SubscribeRepo: subscribeMongoRepo,
		ChatService:   chatServiceClient,
		Log:           log,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
			middleware_metrics.MetricsUnaryInterceptor(),
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewServerServer(resourcesShutdownCtx, server.Deps{
		ServerUsecase: serverUsecase,
		Log:           log,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterServerServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")
	resourcesShutdownCtxCancel()
	log.Print("wait shutdown resources")
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
