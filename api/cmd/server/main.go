package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/aoshimash/rx/api/internal/config"
	"github.com/aoshimash/rx/api/internal/handler"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/aoshimash/rx/api/internal/storage"
	s3storage "github.com/aoshimash/rx/api/internal/storage/s3"
	"github.com/aoshimash/rx/api/internal/store/memory"
	postgresstore "github.com/aoshimash/rx/api/internal/store/postgres"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	var programRepo repository.ProgramRepository
	var logRepo repository.LogRepository
	var planRepo repository.PlanRepository
	var fieldGroupRepo repository.FieldGroupRepository

	ctx := context.Background()

	if cfg.Database.StorageType == "postgres" {
		db, err := postgresstore.NewDB(ctx, cfg.Database)
		if err != nil {
			slog.Error("Failed to initialize PostgreSQL connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		slog.Info("Using PostgreSQL storage backend")
		programRepo = postgresstore.NewProgramRepository(db.Pool())
		logRepo = postgresstore.NewLogRepository(db.Pool())
		planRepo = postgresstore.NewPlanRepository(db.Pool())
		fieldGroupRepo = postgresstore.NewFieldGroupRepository(db.Pool())
	} else {
		slog.Info("Using in-memory storage backend")
		programRepo = memory.NewProgramRepository()
		logRepo = memory.NewLogRepository()
		planRepo = memory.NewPlanRepository()
		fieldGroupRepo = memory.NewFieldGroupRepository()
	}

	var storageProvider storage.Provider
	if cfg.Storage.IsStorageEnabled() {
		switch cfg.Storage.Provider {
		case "s3", "r2":
			if cfg.Storage.S3Config.Bucket == "" {
				slog.Error("STORAGE_BUCKET is required when STORAGE_PROVIDER is set", "provider", cfg.Storage.Provider)
				os.Exit(1)
			}
			var err error
			storageProvider, err = s3storage.New(context.Background(), cfg.Storage.S3Config)
			if err != nil {
				slog.Error("Failed to initialize storage provider", "error", err)
				os.Exit(1)
			}
			slog.Info("Storage provider initialized", "provider", cfg.Storage.Provider)
		default:
			slog.Error("Unsupported storage provider", "provider", cfg.Storage.Provider)
			os.Exit(1)
		}
	} else {
		slog.Info("Video storage is disabled (STORAGE_PROVIDER not set)")
	}

	authProvider := newAuthProvider(cfg.AuthProvider)

	// --- gRPC server ---
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor(authProvider)),
	)

	pb.RegisterHealthServiceServer(grpcServer, handler.NewHealthServer(logRepo))
	pb.RegisterFieldGroupServiceServer(grpcServer, handler.NewFieldGroupServer(fieldGroupRepo))
	pb.RegisterProgramServiceServer(grpcServer, handler.NewProgramServer(programRepo))
	pb.RegisterPlanServiceServer(grpcServer, handler.NewPlanServer(planRepo, programRepo))
	pb.RegisterLogServiceServer(grpcServer, handler.NewLogServer(logRepo, programRepo))
	pb.RegisterVideoServiceServer(grpcServer, handler.NewVideoServer(storageProvider, logger))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		slog.Error("Failed to listen on gRPC port", "port", cfg.GRPCPort, "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	// --- gRPC-Gateway (HTTP) ---
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.EqualFold(key, "authorization") {
				return "authorization", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%s", cfg.GRPCPort)

	registrations := []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
		pb.RegisterHealthServiceHandlerFromEndpoint,
		pb.RegisterFieldGroupServiceHandlerFromEndpoint,
		pb.RegisterProgramServiceHandlerFromEndpoint,
		pb.RegisterPlanServiceHandlerFromEndpoint,
		pb.RegisterLogServiceHandlerFromEndpoint,
		pb.RegisterVideoServiceHandlerFromEndpoint,
	}
	for _, register := range registrations {
		if err := register(ctx, mux, endpoint, opts); err != nil {
			slog.Error("Failed to register gateway handler", "error", err)
			os.Exit(1)
		}
	}

	// Wrap with CORS middleware for the HTTP gateway
	httpHandler := middleware.CORSMiddleware(middleware.DefaultCORSConfig())(mux)

	slog.Info("HTTP server listening", "port", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), httpHandler); err != nil { //nolint:gosec // bind address is configurable via env
		slog.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}

func newAuthProvider(providerName string) middleware.GRPCAuthProvider {
	switch providerName {
	case "stub":
		slog.Info("Using stub authentication provider")
		return &middleware.GRPCStubProvider{}
	default:
		slog.Warn("Unknown auth provider, defaulting to stub", "provider", providerName)
		return &middleware.GRPCStubProvider{}
	}
}
