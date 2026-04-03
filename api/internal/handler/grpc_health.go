package handler

import (
	"context"
	"time"

	"github.com/aoshimash/rx/api/internal/repository"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HealthServer implements the HealthService gRPC server.
type HealthServer struct {
	pb.UnimplementedHealthServiceServer
	logRepo repository.LogRepository
}

// NewHealthServer creates a new HealthServer.
func NewHealthServer(logRepo repository.LogRepository) *HealthServer {
	return &HealthServer{logRepo: logRepo}
}

// Check verifies database connectivity.
func (s *HealthServer) Check(ctx context.Context, _ *pb.CheckRequest) (*pb.CheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, _, _, err := s.logRepo.List(ctx, nil, 1, "")
	if err != nil {
		return nil, status.Error(codes.Unavailable, "unhealthy")
	}
	return &pb.CheckResponse{Status: "healthy"}, nil
}
