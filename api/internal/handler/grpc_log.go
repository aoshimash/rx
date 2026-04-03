package handler

import (
	"context"
	"errors"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogServer implements the LogService gRPC server.
type LogServer struct {
	pb.UnimplementedLogServiceServer
	repo        repository.LogRepository
	programRepo repository.ProgramRepository
}

// NewLogServer creates a new LogServer.
func NewLogServer(repo repository.LogRepository, programRepo repository.ProgramRepository) *LogServer {
	return &LogServer{
		repo:        repo,
		programRepo: programRepo,
	}
}

func logFromProto(req *pb.CreateLogRequest) (*domain.Log, error) {
	if req.GetPerformedAt() == nil {
		return nil, status.Error(codes.InvalidArgument, "performed_at is required")
	}
	log := &domain.Log{
		SessionName:  optionalString(req.GetSessionName()),
		PerformedAt:  timestampToTime(req.GetPerformedAt()),
		Notes:        optionalString(req.GetNotes()),
		Metadata:     structToRawJSON(req.GetMetadata()),
		PlanSnapshot: structToRawJSON(req.GetPlanSnapshot()),
		StartedAt:    timestampToOptionalTime(req.GetStartedAt()),
		FinishedAt:   timestampToOptionalTime(req.GetFinishedAt()),
		ProgramID:    optionalUUID(req.GetProgramId()),
	}

	entries := make([]domain.LogEntry, len(req.GetEntries()))
	for i, entryReq := range req.GetEntries() {
		entry := domain.LogEntry{
			ExerciseName: entryReq.GetExerciseName(),
			Order:        i,
			Fields:       structToMap(entryReq.GetFields()),
			Notes:        optionalString(entryReq.GetNotes()),
			StartedAt:    timestampToOptionalTime(entryReq.GetStartedAt()),
			FinishedAt:   timestampToOptionalTime(entryReq.GetFinishedAt()),
		}

		if len(entryReq.GetSets()) > 0 {
			sets := make([]domain.LogSet, len(entryReq.GetSets()))
			for j, setReq := range entryReq.GetSets() {
				sets[j] = domain.LogSet{
					SetNumber:      int(setReq.GetSetNumber()),
					Fields:         structToMap(setReq.GetFields()),
					VideoObjectKey: optionalString(setReq.GetVideoObjectKey()),
					Notes:          optionalString(setReq.GetNotes()),
				}
			}
			entry.Sets = sets
		}

		entries[i] = entry
	}
	log.Entries = entries

	return log, nil
}

func logFromUpdateProto(req *pb.UpdateLogRequest) (*domain.Log, error) {
	if req.GetPerformedAt() == nil {
		return nil, status.Error(codes.InvalidArgument, "performed_at is required")
	}
	log := &domain.Log{
		SessionName:  optionalString(req.GetSessionName()),
		PerformedAt:  timestampToTime(req.GetPerformedAt()),
		Notes:        optionalString(req.GetNotes()),
		Metadata:     structToRawJSON(req.GetMetadata()),
		PlanSnapshot: structToRawJSON(req.GetPlanSnapshot()),
		StartedAt:    timestampToOptionalTime(req.GetStartedAt()),
		FinishedAt:   timestampToOptionalTime(req.GetFinishedAt()),
		ProgramID:    optionalUUID(req.GetProgramId()),
	}

	entries := make([]domain.LogEntry, len(req.GetEntries()))
	for i, entryReq := range req.GetEntries() {
		entry := domain.LogEntry{
			ExerciseName: entryReq.GetExerciseName(),
			Order:        i,
			Fields:       structToMap(entryReq.GetFields()),
			Notes:        optionalString(entryReq.GetNotes()),
			StartedAt:    timestampToOptionalTime(entryReq.GetStartedAt()),
			FinishedAt:   timestampToOptionalTime(entryReq.GetFinishedAt()),
		}

		if len(entryReq.GetSets()) > 0 {
			sets := make([]domain.LogSet, len(entryReq.GetSets()))
			for j, setReq := range entryReq.GetSets() {
				sets[j] = domain.LogSet{
					SetNumber:      int(setReq.GetSetNumber()),
					Fields:         structToMap(setReq.GetFields()),
					VideoObjectKey: optionalString(setReq.GetVideoObjectKey()),
					Notes:          optionalString(setReq.GetNotes()),
				}
			}
			entry.Sets = sets
		}

		entries[i] = entry
	}
	log.Entries = entries

	return log, nil
}

func (s *LogServer) CreateLog(ctx context.Context, req *pb.CreateLogRequest) (*pb.CreateLogResponse, error) {
	log, err := logFromProto(req)
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := domain.ValidateLog(log); err != nil {
		return nil, domainErrToStatus(err)
	}

	if log.ProgramID != nil {
		if _, err := s.programRepo.GetByID(ctx, *log.ProgramID); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, status.Errorf(codes.InvalidArgument, "program not found: %s", log.ProgramID.String())
			}
			return nil, status.Error(codes.Internal, "failed to retrieve program")
		}
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, status.Error(codes.Internal, "failed to create log")
	}

	return &pb.CreateLogResponse{Log: logToProto(log)}, nil
}

func (s *LogServer) GetLog(ctx context.Context, req *pb.GetLogRequest) (*pb.GetLogResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "log not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve log")
	}

	return &pb.GetLogResponse{Log: logToProto(log)}, nil
}

func (s *LogServer) ListLogs(ctx context.Context, req *pb.ListLogsRequest) (*pb.ListLogsResponse, error) {
	limit := int(req.GetLimit())
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	var programID *uuid.UUID
	if req.GetProgramId() != "" {
		pid, err := uuid.Parse(req.GetProgramId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid program_id format")
		}
		programID = &pid
	}

	performedAtFrom := timestampToOptionalTime(req.GetPerformedAtFrom())
	performedAtTo := timestampToOptionalTime(req.GetPerformedAtTo())

	var logs []*domain.Log
	var nextCursor string
	var hasMore bool
	var listErr error

	if performedAtFrom != nil || performedAtTo != nil {
		logs, nextCursor, hasMore, listErr = s.repo.ListByPerformedAtRange(ctx, programID, performedAtFrom, performedAtTo, limit, req.GetAfter())
	} else {
		logs, nextCursor, hasMore, listErr = s.repo.List(ctx, programID, limit, req.GetAfter())
	}
	if listErr != nil {
		return nil, status.Error(codes.Internal, "failed to list logs")
	}

	data := make([]*pb.Log, len(logs))
	for i, l := range logs {
		data[i] = logToProto(l)
	}

	return &pb.ListLogsResponse{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *LogServer) UpdateLog(ctx context.Context, req *pb.UpdateLogRequest) (*pb.UpdateLogResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "log not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve log")
	}

	updated, err := logFromUpdateProto(req)
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	if err := domain.ValidateLog(updated); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.repo.Update(ctx, updated); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "log not found")
		}
		return nil, status.Error(codes.Internal, "failed to update log")
	}

	return &pb.UpdateLogResponse{Log: logToProto(updated)}, nil
}

func (s *LogServer) DeleteLog(ctx context.Context, req *pb.DeleteLogRequest) (*pb.DeleteLogResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "log not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve log")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "log not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete log")
	}

	return &pb.DeleteLogResponse{}, nil
}
