package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/storage"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VideoServer implements the VideoService gRPC server.
type VideoServer struct {
	pb.UnimplementedVideoServiceServer
	storage storage.Provider
	logger  *slog.Logger
}

// NewVideoServer creates a new VideoServer.
func NewVideoServer(storage storage.Provider, logger *slog.Logger) *VideoServer {
	return &VideoServer{
		storage: storage,
		logger:  logger,
	}
}

func (s *VideoServer) GenerateUploadURL(ctx context.Context, req *pb.GenerateUploadURLRequest) (*pb.GenerateUploadURLResponse, error) {
	if s.storage == nil {
		return nil, status.Error(codes.Unavailable, "video storage is not configured")
	}

	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	if req.GetContentType() == "" {
		return nil, status.Error(codes.InvalidArgument, "content_type is required")
	}
	if req.GetFilename() == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is required")
	}
	if req.GetContentLength() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "content_length must be positive")
	}

	uploadReq := storage.UploadURLRequest{
		ContentType:   req.GetContentType(),
		Filename:      req.GetFilename(),
		UserID:        userID,
		ContentLength: req.GetContentLength(),
	}

	resp, err := s.storage.GenerateUploadURL(ctx, uploadReq)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidContentType) {
			return nil, status.Error(codes.InvalidArgument, "content_type must be video/*")
		}
		if errors.Is(err, storage.ErrFileTooLarge) {
			return nil, status.Error(codes.InvalidArgument, "file size exceeds maximum allowed")
		}
		s.logger.Error("failed to generate upload URL", "error", err)
		return nil, status.Error(codes.Internal, "failed to generate upload URL")
	}

	return &pb.GenerateUploadURLResponse{
		UploadUrl: resp.UploadURL,
		ObjectKey: resp.ObjectKey,
		ExpiresIn: int32(resp.ExpiresIn.Seconds()),
	}, nil
}

func (s *VideoServer) GenerateDownloadURL(ctx context.Context, req *pb.GenerateDownloadURLRequest) (*pb.GenerateDownloadURLResponse, error) {
	if s.storage == nil {
		return nil, status.Error(codes.Unavailable, "video storage is not configured")
	}

	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	if req.GetObjectKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "object_key is required")
	}

	downloadReq := storage.DownloadURLRequest{
		ObjectKey: req.GetObjectKey(),
		UserID:    userID,
	}

	resp, err := s.storage.GenerateDownloadURL(ctx, downloadReq)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidObjectKey) {
			return nil, status.Error(codes.InvalidArgument, "invalid or unauthorized object key")
		}
		s.logger.Error("failed to generate download URL", "error", err)
		return nil, status.Error(codes.Internal, "failed to generate download URL")
	}

	return &pb.GenerateDownloadURLResponse{
		DownloadUrl: resp.DownloadURL,
		ExpiresIn:   int32(resp.ExpiresIn.Seconds()),
	}, nil
}
