package handler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	pb "github.com/aoshimash/rx/api/pkg/gen/rx/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- Error conversion ---

func domainErrToStatus(err error) error {
	if err == nil {
		return nil
	}
	if ve, ok := err.(*domain.ValidationError); ok {
		return status.Errorf(codes.InvalidArgument, "%s: %s", ve.Field, ve.Message)
	}
	if de, ok := err.(*domain.DomainError); ok {
		switch de.Code {
		case domain.ErrorCodeNotFound:
			return status.Error(codes.NotFound, de.Message)
		case domain.ErrorCodeConflict:
			return status.Error(codes.AlreadyExists, de.Message)
		default:
			return status.Error(codes.Internal, de.Message)
		}
	}
	if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, "resource not found")
	}
	return status.Error(codes.Internal, "internal error")
}

// --- UUID helpers ---

func parseUUID(s, field string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "%s is required", field)
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "invalid %s format: %s", field, s)
	}
	return id, nil
}

func optionalUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func uuidToString(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

func optionalUUIDToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

// --- Timestamp helpers ---

func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func optionalTimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func timestampToOptionalTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// --- Struct <-> map helpers ---

func mapToStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return s
}

func structToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

func rawJSONToStruct(raw json.RawMessage) *structpb.Struct {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	return mapToStruct(m)
}

func structToRawJSON(s *structpb.Struct) json.RawMessage {
	if s == nil {
		return nil
	}
	m := s.AsMap()
	if m == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

// --- String <-> *string helpers ---

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// --- DateOnly helpers ---

func dateOnlyToString(d *domain.DateOnly) string {
	if d == nil {
		return ""
	}
	return time.Time(*d).Format("2006-01-02")
}

func stringToDateOnly(s string) *domain.DateOnly {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	d := domain.DateOnly(t)
	return &d
}

// --- FieldDef conversion ---

func fieldDefToProto(fd domain.FieldDef) *pb.FieldDef {
	return &pb.FieldDef{
		Name:        fd.Name,
		Type:        fieldTypeToProto(fd.Type),
		Options:     fd.Options,
		Description: fd.Description,
	}
}

func fieldDefFromProto(fd *pb.FieldDef) domain.FieldDef {
	return domain.FieldDef{
		Name:        fd.GetName(),
		Type:        fieldTypeFromProto(fd.GetType()),
		Options:     fd.GetOptions(),
		Description: fd.GetDescription(),
	}
}

func fieldTypeToProto(t string) pb.FieldType {
	switch t {
	case "text":
		return pb.FieldType_FIELD_TYPE_TEXT
	case "number":
		return pb.FieldType_FIELD_TYPE_NUMBER
	case "select":
		return pb.FieldType_FIELD_TYPE_SELECT
	case "video":
		return pb.FieldType_FIELD_TYPE_VIDEO
	default:
		return pb.FieldType_FIELD_TYPE_UNSPECIFIED
	}
}

func fieldTypeFromProto(t pb.FieldType) string {
	switch t {
	case pb.FieldType_FIELD_TYPE_TEXT:
		return "text"
	case pb.FieldType_FIELD_TYPE_NUMBER:
		return "number"
	case pb.FieldType_FIELD_TYPE_SELECT:
		return "select"
	case pb.FieldType_FIELD_TYPE_VIDEO:
		return "video"
	default:
		return ""
	}
}

func fieldDefsToProto(defs []domain.FieldDef) []*pb.FieldDef {
	if defs == nil {
		return nil
	}
	result := make([]*pb.FieldDef, len(defs))
	for i, fd := range defs {
		result[i] = fieldDefToProto(fd)
	}
	return result
}

func fieldDefsFromProto(defs []*pb.FieldDef) []domain.FieldDef {
	if defs == nil {
		return nil
	}
	result := make([]domain.FieldDef, len(defs))
	for i, fd := range defs {
		result[i] = fieldDefFromProto(fd)
	}
	return result
}

// --- FieldGroup conversion ---

func fieldGroupToProto(fg *domain.FieldGroup) *pb.FieldGroup {
	return &pb.FieldGroup{
		Id:            uuidToString(fg.ID),
		Name:          fg.Name,
		Description:   derefString(fg.Description),
		ProgramFields: fieldDefsToProto(fg.ProgramFields),
		LogFields:     fieldDefsToProto(fg.LogFields),
		CreatedAt:     timeToTimestamp(fg.CreatedAt),
		UpdatedAt:     timeToTimestamp(fg.UpdatedAt),
	}
}

// --- Program conversion ---

func programGroupToProto(g domain.ProgramGroup) *pb.ProgramGroup {
	return &pb.ProgramGroup{
		Id:            uuidToString(g.ID),
		ProgramId:     uuidToString(g.ProgramID),
		ParentGroupId: optionalUUIDToString(g.ParentGroupID),
		Name:          g.Name,
		Order:         int32(g.Order),
		Notes:         derefString(g.Notes),
	}
}

func programSessionEntryToProto(e domain.ProgramSessionEntry) *pb.ProgramSessionEntry {
	return &pb.ProgramSessionEntry{
		Id:           uuidToString(e.ID),
		SessionId:    uuidToString(e.SessionID),
		Order:        int32(e.Order),
		ExerciseName: e.ExerciseName,
		Fields:       mapToStruct(e.Fields),
		Notes:        derefString(e.Notes),
	}
}

func programSessionToProto(s domain.ProgramSession) *pb.ProgramSession {
	entries := make([]*pb.ProgramSessionEntry, len(s.Entries))
	for i, e := range s.Entries {
		entries[i] = programSessionEntryToProto(e)
	}
	return &pb.ProgramSession{
		Id:           uuidToString(s.ID),
		ProgramId:    uuidToString(s.ProgramID),
		GroupId:      optionalUUIDToString(s.GroupID),
		SessionName:  s.SessionName,
		Order:        int32(s.Order),
		FieldGroupId: optionalUUIDToString(s.FieldGroupID),
		Date:         dateOnlyToString(s.Date),
		Entries:      entries,
	}
}

func programToProto(p *domain.Program) *pb.Program {
	groups := make([]*pb.ProgramGroup, len(p.Groups))
	for i, g := range p.Groups {
		groups[i] = programGroupToProto(g)
	}
	sessions := make([]*pb.ProgramSession, len(p.Sessions))
	for i, s := range p.Sessions {
		sessions[i] = programSessionToProto(s)
	}
	return &pb.Program{
		Id:        uuidToString(p.ID),
		Name:      p.Name,
		Notes:     derefString(p.Notes),
		Groups:    groups,
		Sessions:  sessions,
		CreatedAt: timeToTimestamp(p.CreatedAt),
		UpdatedAt: timeToTimestamp(p.UpdatedAt),
	}
}

// --- Plan conversion ---

func planSessionEntryToProto(e domain.PlanSessionEntry) *pb.PlanSessionEntry {
	return &pb.PlanSessionEntry{
		Id:           uuidToString(e.ID),
		SessionId:    uuidToString(e.SessionID),
		Order:        int32(e.Order),
		ExerciseName: e.ExerciseName,
		Fields:       mapToStruct(e.Fields),
		Notes:        derefString(e.Notes),
	}
}

func planSessionToProto(s domain.PlanSession) *pb.PlanSession {
	entries := make([]*pb.PlanSessionEntry, len(s.Entries))
	for i, e := range s.Entries {
		entries[i] = planSessionEntryToProto(e)
	}
	return &pb.PlanSession{
		Id:              uuidToString(s.ID),
		PlanId:          uuidToString(s.PlanID),
		SessionName:     s.SessionName,
		Order:           int32(s.Order),
		FieldGroupId:    optionalUUIDToString(s.FieldGroupID),
		Date:            dateOnlyToString(s.Date),
		SourceProgramId: optionalUUIDToString(s.SourceProgramID),
		SourceSessionId: optionalUUIDToString(s.SourceSessionID),
		Entries:         entries,
	}
}

func planToProto(p *domain.Plan) *pb.Plan {
	sessions := make([]*pb.PlanSession, len(p.Sessions))
	for i, s := range p.Sessions {
		sessions[i] = planSessionToProto(s)
	}
	return &pb.Plan{
		Id:        uuidToString(p.ID),
		Name:      derefString(p.Name),
		Notes:     derefString(p.Notes),
		Sessions:  sessions,
		CreatedAt: timeToTimestamp(p.CreatedAt),
		UpdatedAt: timeToTimestamp(p.UpdatedAt),
	}
}

// --- Log conversion ---

func logSetToProto(s domain.LogSet) *pb.LogSet {
	return &pb.LogSet{
		Id:             uuidToString(s.ID),
		EntryId:        uuidToString(s.EntryID),
		SetNumber:      int32(s.SetNumber),
		Fields:         mapToStruct(s.Fields),
		VideoObjectKey: derefString(s.VideoObjectKey),
		Notes:          derefString(s.Notes),
	}
}

func logEntryToProto(e domain.LogEntry) *pb.LogEntry {
	sets := make([]*pb.LogSet, len(e.Sets))
	for i, s := range e.Sets {
		sets[i] = logSetToProto(s)
	}
	return &pb.LogEntry{
		Id:           uuidToString(e.ID),
		LogId:        uuidToString(e.LogID),
		Order:        int32(e.Order),
		ExerciseName: e.ExerciseName,
		Fields:       mapToStruct(e.Fields),
		Sets:         sets,
		Notes:        derefString(e.Notes),
		StartedAt:    optionalTimeToTimestamp(e.StartedAt),
		FinishedAt:   optionalTimeToTimestamp(e.FinishedAt),
	}
}

func logToProto(l *domain.Log) *pb.Log {
	entries := make([]*pb.LogEntry, len(l.Entries))
	for i, e := range l.Entries {
		entries[i] = logEntryToProto(e)
	}
	return &pb.Log{
		Id:           uuidToString(l.ID),
		ProgramId:    optionalUUIDToString(l.ProgramID),
		SessionName:  derefString(l.SessionName),
		PerformedAt:  timeToTimestamp(l.PerformedAt),
		Notes:        derefString(l.Notes),
		Metadata:     rawJSONToStruct(l.Metadata),
		PlanSnapshot: rawJSONToStruct(l.PlanSnapshot),
		StartedAt:    optionalTimeToTimestamp(l.StartedAt),
		FinishedAt:   optionalTimeToTimestamp(l.FinishedAt),
		Entries:      entries,
		CreatedAt:    timeToTimestamp(l.CreatedAt),
		UpdatedAt:    timeToTimestamp(l.UpdatedAt),
	}
}

// --- Plan session from proto ---

func planSessionFromProto(req *pb.PlanSessionCreate) (domain.PlanSession, error) {
	sess := domain.PlanSession{
		SessionName:     req.GetSessionName(),
		Order:           int(req.GetOrder()),
		FieldGroupID:    optionalUUID(req.GetFieldGroupId()),
		Date:            stringToDateOnly(req.GetDate()),
		SourceProgramID: optionalUUID(req.GetSourceProgramId()),
		SourceSessionID: optionalUUID(req.GetSourceSessionId()),
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
	sess.Entries = entries

	return sess, nil
}

func planSessionsFromProto(reqs []*pb.PlanSessionCreate) ([]domain.PlanSession, error) {
	result := make([]domain.PlanSession, len(reqs))
	for i, req := range reqs {
		sess, err := planSessionFromProto(req)
		if err != nil {
			return nil, err
		}
		result[i] = sess
	}
	return result, nil
}
