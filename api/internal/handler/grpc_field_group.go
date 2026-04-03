package handler

import (
	"context"
	"errors"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FieldGroupServer implements the FieldGroupService gRPC server.
type FieldGroupServer struct {
	pb.UnimplementedFieldGroupServiceServer
	repo repository.FieldGroupRepository
}

// NewFieldGroupServer creates a new FieldGroupServer.
func NewFieldGroupServer(repo repository.FieldGroupRepository) *FieldGroupServer {
	return &FieldGroupServer{repo: repo}
}

func (s *FieldGroupServer) ListFieldGroups(ctx context.Context, _ *pb.ListFieldGroupsRequest) (*pb.ListFieldGroupsResponse, error) {
	userID := middleware.GetUserID(ctx)

	groups, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list field groups")
	}

	data := make([]*pb.FieldGroup, len(groups))
	for i := range groups {
		data[i] = fieldGroupToProto(&groups[i])
	}
	return &pb.ListFieldGroupsResponse{Data: data}, nil
}

func (s *FieldGroupServer) CreateFieldGroup(ctx context.Context, req *pb.CreateFieldGroupRequest) (*pb.CreateFieldGroupResponse, error) {
	userID := middleware.GetUserID(ctx)

	fg := &domain.FieldGroup{
		Name:          req.GetName(),
		Description:   optionalString(req.GetDescription()),
		ProgramFields: fieldDefsFromProto(req.GetProgramFields()),
		LogFields:     fieldDefsFromProto(req.GetLogFields()),
	}

	if err := domain.ValidateFieldGroup(fg); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.repo.Create(ctx, userID, fg); err != nil {
		if de, ok := err.(*domain.DomainError); ok && de.Code == domain.ErrorCodeConflict {
			return nil, status.Error(codes.AlreadyExists, "a field group with this name already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create field group")
	}

	return &pb.CreateFieldGroupResponse{FieldGroup: fieldGroupToProto(fg)}, nil
}

func (s *FieldGroupServer) GetFieldGroup(ctx context.Context, req *pb.GetFieldGroupRequest) (*pb.GetFieldGroupResponse, error) {
	userID := middleware.GetUserID(ctx)

	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	fg, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "field group not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve field group")
	}

	return &pb.GetFieldGroupResponse{FieldGroup: fieldGroupToProto(fg)}, nil
}

func (s *FieldGroupServer) UpdateFieldGroup(ctx context.Context, req *pb.UpdateFieldGroupRequest) (*pb.UpdateFieldGroupResponse, error) {
	userID := middleware.GetUserID(ctx)

	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "field group not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve field group")
	}

	existing.Name = req.GetName()
	existing.Description = optionalString(req.GetDescription())
	existing.ProgramFields = fieldDefsFromProto(req.GetProgramFields())
	existing.LogFields = fieldDefsFromProto(req.GetLogFields())

	if err := domain.ValidateFieldGroup(existing); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.repo.Update(ctx, userID, existing); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "field group not found")
		}
		if de, ok := err.(*domain.DomainError); ok && de.Code == domain.ErrorCodeConflict {
			return nil, status.Error(codes.AlreadyExists, "a field group with this name already exists")
		}
		return nil, status.Error(codes.Internal, "failed to update field group")
	}

	return &pb.UpdateFieldGroupResponse{FieldGroup: fieldGroupToProto(existing)}, nil
}

func (s *FieldGroupServer) DeleteFieldGroup(ctx context.Context, req *pb.DeleteFieldGroupRequest) (*pb.DeleteFieldGroupResponse, error) {
	userID := middleware.GetUserID(ctx)

	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, userID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "field group not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete field group")
	}

	return &pb.DeleteFieldGroupResponse{}, nil
}
