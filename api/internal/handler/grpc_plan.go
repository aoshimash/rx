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

// PlanServer implements the PlanService gRPC server.
type PlanServer struct {
	pb.UnimplementedPlanServiceServer
	planRepo    repository.PlanRepository
	programRepo repository.ProgramRepository
}

// NewPlanServer creates a new PlanServer.
func NewPlanServer(planRepo repository.PlanRepository, programRepo repository.ProgramRepository) *PlanServer {
	return &PlanServer{
		planRepo:    planRepo,
		programRepo: programRepo,
	}
}

func (s *PlanServer) GetPlan(ctx context.Context, _ *pb.GetPlanRequest) (*pb.GetPlanResponse, error) {
	userID := middleware.GetUserID(ctx)

	plan, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve plan")
	}

	return &pb.GetPlanResponse{Plan: planToProto(plan)}, nil
}

func (s *PlanServer) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.CreatePlanResponse, error) {
	userID := middleware.GetUserID(ctx)

	plan := &domain.Plan{
		Name:  optionalString(req.GetName()),
		Notes: optionalString(req.GetNotes()),
	}

	if len(req.GetSessions()) > 0 {
		sessions, err := planSessionsFromProto(req.GetSessions())
		if err != nil {
			return nil, domainErrToStatus(err)
		}
		plan.Sessions = sessions
	}

	if err := domain.ValidatePlan(plan); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.planRepo.Create(ctx, userID, plan); err != nil {
		if de, ok := err.(*domain.DomainError); ok && de.Code == domain.ErrorCodeConflict {
			return nil, status.Error(codes.AlreadyExists, "a plan already exists for this user")
		}
		return nil, status.Error(codes.Internal, "failed to create plan")
	}

	created, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve created plan")
	}

	return &pb.CreatePlanResponse{Plan: planToProto(created)}, nil
}

func (s *PlanServer) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.UpdatePlanResponse, error) {
	userID := middleware.GetUserID(ctx)

	existing, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve plan")
	}

	plan := &domain.Plan{
		Name:  optionalString(req.GetName()),
		Notes: optionalString(req.GetNotes()),
	}

	if req.GetSessions() != nil {
		sessions, err := planSessionsFromProto(req.GetSessions())
		if err != nil {
			return nil, domainErrToStatus(err)
		}
		plan.Sessions = sessions
	} else {
		plan.Sessions = existing.Sessions
	}

	if err := domain.ValidatePlan(plan); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.planRepo.Update(ctx, userID, plan); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to update plan")
	}

	updated, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve updated plan")
	}

	return &pb.UpdatePlanResponse{Plan: planToProto(updated)}, nil
}

func (s *PlanServer) DeletePlan(ctx context.Context, _ *pb.DeletePlanRequest) (*pb.DeletePlanResponse, error) {
	userID := middleware.GetUserID(ctx)

	if err := s.planRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete plan")
	}

	return &pb.DeletePlanResponse{}, nil
}

func (s *PlanServer) AddPlanSessions(ctx context.Context, req *pb.AddPlanSessionsRequest) (*pb.AddPlanSessionsResponse, error) {
	userID := middleware.GetUserID(ctx)

	if len(req.GetSessions()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one session is required")
	}

	sessions, err := planSessionsFromProto(req.GetSessions())
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	for i := range sessions {
		if err := domain.ValidatePlanSession(&sessions[i]); err != nil {
			return nil, domainErrToStatus(err)
		}
	}

	if err := s.planRepo.AddSessions(ctx, userID, sessions); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, "failed to add sessions")
	}

	plan, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve updated plan")
	}

	return &pb.AddPlanSessionsResponse{Plan: planToProto(plan)}, nil
}

func (s *PlanServer) UpdatePlanSession(ctx context.Context, req *pb.UpdatePlanSessionRequest) (*pb.UpdatePlanSessionResponse, error) {
	userID := middleware.GetUserID(ctx)

	sessionID, err := parseUUID(req.GetSessionId(), "session_id")
	if err != nil {
		return nil, err
	}

	entries := make([]domain.PlanSessionEntry, len(req.GetEntries()))
	for j, e := range req.GetEntries() {
		entries[j] = domain.PlanSessionEntry{
			ExerciseName: e.GetExerciseName(),
			Order:        int(e.GetOrder()),
			Fields:       structToMap(e.GetFields()),
			Notes:        optionalString(e.GetNotes()),
		}
	}

	sess := &domain.PlanSession{
		ID:              sessionID,
		SessionName:     req.GetSessionName(),
		Order:           int(req.GetOrder()),
		FieldGroupID:    optionalUUID(req.GetFieldGroupId()),
		Date:            stringToDateOnly(req.GetDate()),
		SourceProgramID: optionalUUID(req.GetSourceProgramId()),
		SourceSessionID: optionalUUID(req.GetSourceSessionId()),
		Entries:         entries,
	}

	if err := domain.ValidatePlanSession(sess); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.planRepo.UpdateSession(ctx, userID, sess); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan or session not found")
		}
		return nil, status.Error(codes.Internal, "failed to update session")
	}

	return &pb.UpdatePlanSessionResponse{Session: planSessionToProto(*sess)}, nil
}

func (s *PlanServer) DeletePlanSession(ctx context.Context, req *pb.DeletePlanSessionRequest) (*pb.DeletePlanSessionResponse, error) {
	userID := middleware.GetUserID(ctx)

	if _, err := parseUUID(req.GetSessionId(), "session_id"); err != nil {
		return nil, err
	}

	if err := s.planRepo.DeleteSession(ctx, userID, req.GetSessionId()); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan or session not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete session")
	}

	return &pb.DeletePlanSessionResponse{}, nil
}

func (s *PlanServer) ExpandProgram(ctx context.Context, req *pb.ExpandProgramRequest) (*pb.ExpandProgramResponse, error) {
	userID := middleware.GetUserID(ctx)

	programID, err := parseUUID(req.GetProgramId(), "program_id")
	if err != nil {
		return nil, err
	}

	program, err := s.programRepo.GetByID(ctx, programID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve program")
	}

	planSessions := make([]domain.PlanSession, len(program.Sessions))
	for i, ps := range program.Sessions {
		entries := make([]domain.PlanSessionEntry, len(ps.Entries))
		for j, e := range ps.Entries {
			entries[j] = domain.PlanSessionEntry{
				Order:        e.Order,
				ExerciseName: e.ExerciseName,
				Fields:       e.Fields,
				Notes:        e.Notes,
			}
		}

		srcProgramID := program.ID
		srcSessionID := ps.ID
		planSessions[i] = domain.PlanSession{
			SessionName:     ps.SessionName,
			Order:           ps.Order,
			FieldGroupID:    ps.FieldGroupID,
			Date:            ps.Date,
			SourceProgramID: &srcProgramID,
			SourceSessionID: &srcSessionID,
			Entries:         entries,
		}
	}

	_, err = s.planRepo.GetByUserID(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		plan := &domain.Plan{Sessions: planSessions}
		if err := s.planRepo.Create(ctx, userID, plan); err != nil {
			return nil, status.Error(codes.Internal, "failed to create plan")
		}
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to check existing plan")
	} else {
		if err := s.planRepo.AddSessions(ctx, userID, planSessions); err != nil {
			return nil, status.Error(codes.Internal, "failed to add sessions to plan")
		}
	}

	plan, err := s.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve plan")
	}

	return &pb.ExpandProgramResponse{Plan: planToProto(plan)}, nil
}
