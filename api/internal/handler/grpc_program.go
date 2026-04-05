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

// ProgramServer implements the ProgramService gRPC server.
type ProgramServer struct {
	pb.UnimplementedProgramServiceServer
	repo repository.ProgramRepository
}

// NewProgramServer creates a new ProgramServer.
func NewProgramServer(repo repository.ProgramRepository) *ProgramServer {
	return &ProgramServer{repo: repo}
}

// parseGroupsFromProto converts proto ProgramGroupCreate to domain ProgramGroup slice
// with temp_id resolution and topological sorting.
func parseGroupsFromProto(groups []*pb.ProgramGroupCreate) ([]domain.ProgramGroup, map[string]uuid.UUID, error) {
	tempIDMap := make(map[string]uuid.UUID, len(groups))
	result := make([]domain.ProgramGroup, len(groups))

	for i, g := range groups {
		id := uuid.New()
		result[i] = domain.ProgramGroup{
			ID:    id,
			Name:  g.GetName(),
			Order: int(g.GetOrder()),
			Notes: optionalString(g.GetNotes()),
		}
		if g.GetTempId() != "" {
			if _, exists := tempIDMap[g.GetTempId()]; exists {
				return nil, nil, &domain.ValidationError{
					Field:   "groups[].temp_id",
					Message: "duplicate temp_id: " + g.GetTempId(),
				}
			}
			tempIDMap[g.GetTempId()] = id
		}
	}

	for i, g := range groups {
		if g.GetParentGroupId() == "" {
			continue
		}
		parentUUID, ok := tempIDMap[g.GetParentGroupId()]
		if !ok {
			return nil, nil, &domain.ValidationError{
				Field:   "groups[].parent_group_id",
				Message: "references unknown temp_id: " + g.GetParentGroupId(),
			}
		}
		result[i].ParentGroupID = &parentUUID
	}

	sorted, err := topoSortGroups(result)
	if err != nil {
		return nil, nil, err
	}

	return sorted, tempIDMap, nil
}

// parseSessionsFromProto converts proto ProgramSessionCreate to domain ProgramSession slice.
func parseSessionsFromProto(sessions []*pb.ProgramSessionCreate, tempIDMap map[string]uuid.UUID) ([]domain.ProgramSession, error) {
	result := make([]domain.ProgramSession, len(sessions))
	for i, sessReq := range sessions {
		sess := domain.ProgramSession{
			SessionName: sessReq.GetSessionName(),
			Order:       int(sessReq.GetOrder()),
		}

		if sessReq.GetGroupId() != "" {
			gid, ok := tempIDMap[sessReq.GetGroupId()]
			if !ok {
				return nil, &domain.ValidationError{
					Field:   "sessions[].group_id",
					Message: "references unknown temp_id: " + sessReq.GetGroupId(),
				}
			}
			sess.GroupID = &gid
		}

		sess.FieldGroupID = optionalUUID(sessReq.GetFieldGroupId())
		sess.Date = stringToDateOnly(sessReq.GetDate())

		entries := make([]domain.ProgramSessionEntry, len(sessReq.GetEntries()))
		for j, e := range sessReq.GetEntries() {
			entries[j] = domain.ProgramSessionEntry{
				ExerciseName: e.GetExerciseName(),
				Order:        int(e.GetOrder()),
				Fields:       structToMap(e.GetFields()),
				Notes:        optionalString(e.GetNotes()),
			}
		}
		sess.Entries = entries
		result[i] = sess
	}
	return result, nil
}

func (s *ProgramServer) CreateProgram(ctx context.Context, req *pb.CreateProgramRequest) (*pb.CreateProgramResponse, error) {
	exists, err := s.repo.ExistsByName(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check program name")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "a program with this name already exists")
	}

	groups, tempIDMap, err := parseGroupsFromProto(req.GetGroups())
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	sessions, err := parseSessionsFromProto(req.GetSessions(), tempIDMap)
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	program := &domain.Program{
		Name:     req.GetName(),
		Notes:    optionalString(req.GetNotes()),
		Status:   programStatusFromProto(req.GetStatus()),
		Groups:   groups,
		Sessions: sessions,
	}

	if err := domain.ValidateProgram(program); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.repo.Create(ctx, program); err != nil {
		return nil, status.Error(codes.Internal, "failed to create program")
	}

	return &pb.CreateProgramResponse{Program: programToProto(program)}, nil
}

func (s *ProgramServer) GetProgram(ctx context.Context, req *pb.GetProgramRequest) (*pb.GetProgramResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	program, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve program")
	}

	return &pb.GetProgramResponse{Program: programToProto(program)}, nil
}

func (s *ProgramServer) ListPrograms(ctx context.Context, req *pb.ListProgramsRequest) (*pb.ListProgramsResponse, error) {
	limit := int(req.GetLimit())
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	programs, nextCursor, hasMore, err := s.repo.List(ctx, limit, req.GetAfter())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list programs")
	}

	data := make([]*pb.Program, len(programs))
	for i, p := range programs {
		data[i] = programToProto(p)
	}

	return &pb.ListProgramsResponse{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *ProgramServer) UpdateProgram(ctx context.Context, req *pb.UpdateProgramRequest) (*pb.UpdateProgramResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve program")
	}

	if req.GetName() != existing.Name {
		nameExists, err := s.repo.ExistsByName(ctx, req.GetName())
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to check program name")
		}
		if nameExists {
			return nil, status.Error(codes.AlreadyExists, "a program with this name already exists")
		}
	}

	groups, tempIDMap, err := parseGroupsFromProto(req.GetGroups())
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	sessions, err := parseSessionsFromProto(req.GetSessions(), tempIDMap)
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	var reqStatus domain.ProgramStatus
	if req.GetStatus() == pb.ProgramStatus_PROGRAM_STATUS_UNSPECIFIED {
		reqStatus = existing.Status
	} else {
		reqStatus = programStatusFromProto(req.GetStatus())
	}

	program := &domain.Program{
		ID:       existing.ID,
		Name:     req.GetName(),
		Notes:    optionalString(req.GetNotes()),
		Status:   reqStatus,
		Groups:   groups,
		Sessions: sessions,
	}

	if err := domain.ValidateProgram(program); err != nil {
		return nil, domainErrToStatus(err)
	}

	if err := s.repo.Update(ctx, program); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to update program")
	}

	return &pb.UpdateProgramResponse{Program: programToProto(program)}, nil
}

func (s *ProgramServer) UpdateProgramStatus(ctx context.Context, req *pb.UpdateProgramStatusRequest) (*pb.UpdateProgramStatusResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	newStatus := programStatusFromProto(req.GetStatus())
	switch newStatus {
	case domain.ProgramStatusDraft, domain.ProgramStatusPublished:
		// valid
	default:
		return nil, status.Error(codes.InvalidArgument, "status must be draft or published")
	}

	program, err := s.repo.UpdateStatus(ctx, id, newStatus)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to update program status")
	}

	return &pb.UpdateProgramStatusResponse{Program: programToProto(program)}, nil
}

func (s *ProgramServer) DeleteProgram(ctx context.Context, req *pb.DeleteProgramRequest) (*pb.DeleteProgramResponse, error) {
	id, err := parseUUID(req.GetId(), "id")
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve program")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "program not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete program")
	}

	return &pb.DeleteProgramResponse{}, nil
}
