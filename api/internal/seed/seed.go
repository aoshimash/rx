// Package seed provides local development seed data.
// It is intended for use with in-memory storage and local postgres only.
package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

func intPtr(v int) *int             { return &v }
func float64Ptr(v float64) *float64 { return &v }
func stringPtr(v string) *string    { return &v }

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// Run inserts sample powerlifting data into the provided repositories.
//
// Program lifecycle simulated:
//
//	Completed:  SBD Peaking — Jan 2026        (all sessions logged)
//	Completed:  Upper/Lower Split — Feb 2026   (all sessions logged)
//	Ongoing:    Block Periodization — Mar 2026 (3/6 sessions logged)
//	Created:    Competition Prep — Apr 2026    (not started)
func Run(ctx context.Context, programRepo repository.ProgramRepository, logRepo repository.LogRepository) error {
	// ── Program 1: SBD Peaking (completed) ──────────────────────────────────
	prog1 := sbdPeaking()
	if err := programRepo.Create(ctx, prog1); err != nil {
		return fmt.Errorf("create SBD Peaking program: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, prog1.ID, domain.ProgramStatusOngoing); err != nil {
		return fmt.Errorf("update SBD Peaking to ongoing: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, prog1.ID, domain.ProgramStatusCompleted); err != nil {
		return fmt.Errorf("update SBD Peaking to completed: %w", err)
	}
	slog.Info("seeded program", "name", prog1.Name, "status", "completed")

	// Logs for all 4 sessions
	for i, s := range prog1.Sessions {
		log := logForSession(prog1.ID, s, time.Date(2026, 1, 6+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog1.Name, s.SessionName, err)
		}
	}

	// ── Program 2: Upper/Lower Split (completed) ────────────────────────────
	prog2 := upperLowerSplit()
	if err := programRepo.Create(ctx, prog2); err != nil {
		return fmt.Errorf("create Upper/Lower Split program: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, prog2.ID, domain.ProgramStatusOngoing); err != nil {
		return fmt.Errorf("update Upper/Lower to ongoing: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, prog2.ID, domain.ProgramStatusCompleted); err != nil {
		return fmt.Errorf("update Upper/Lower to completed: %w", err)
	}
	slog.Info("seeded program", "name", prog2.Name, "status", "completed")

	for i, s := range prog2.Sessions {
		log := logForSession(prog2.ID, s, time.Date(2026, 2, 2+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog2.Name, s.SessionName, err)
		}
	}

	// ── Program 3: Block Periodization (ongoing, 3/6 logged) ────────────────
	prog3 := blockPeriodization()
	if err := programRepo.Create(ctx, prog3); err != nil {
		return fmt.Errorf("create Block Periodization program: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, prog3.ID, domain.ProgramStatusOngoing); err != nil {
		return fmt.Errorf("update Block Periodization to ongoing: %w", err)
	}
	slog.Info("seeded program", "name", prog3.Name, "status", "ongoing")

	for i := 0; i < 3; i++ {
		s := prog3.Sessions[i]
		log := logForSession(prog3.ID, s, time.Date(2026, 3, 3+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog3.Name, s.SessionName, err)
		}
	}

	// ── Program 4: Competition Prep (created, no logs) ──────────────────────
	prog4 := competitionPrep()
	if err := programRepo.Create(ctx, prog4); err != nil {
		return fmt.Errorf("create Competition Prep program: %w", err)
	}
	slog.Info("seeded program", "name", prog4.Name, "status", "created")

	slog.Info("seed data loaded successfully")
	return nil
}

// ── Program builders ────────────────────────────────────────────────────────

func sbdPeaking() *domain.Program {
	progID := uuid.New()
	s1ID, s2ID, s3ID, s4ID := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	return &domain.Program{
		ID:     progID,
		Name:   "SBD Peaking — Jan 2026",
		Status: domain.ProgramStatusCreated,
		Notes:  stringPtr("4-session peaking block focusing on competition lifts"),
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Heavy Squat", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(5), Reps: intPtr(3), LoadKg: float64Ptr(140), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Pause Squat", Sets: intPtr(3), Reps: intPtr(3), LoadKg: float64Ptr(120), Metadata: mustJSON(map[string]any{"rpe": 7})},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Leg Press", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(200)},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Heavy Bench", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Bench Press", Sets: intPtr(5), Reps: intPtr(3), LoadKg: float64Ptr(100), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Close Grip Bench Press", Sets: intPtr(3), Reps: intPtr(5), LoadKg: float64Ptr(85), Metadata: mustJSON(map[string]any{"rpe": 7})},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Dumbbell Fly", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(16)},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Heavy Deadlift", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Conventional Deadlift", Sets: intPtr(5), Reps: intPtr(2), LoadKg: float64Ptr(180), Metadata: mustJSON(map[string]any{"rpe": 8.5})},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Deficit Deadlift", Sets: intPtr(3), Reps: intPtr(4), LoadKg: float64Ptr(150), Metadata: mustJSON(map[string]any{"rpe": 7})},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(8), LoadKg: float64Ptr(80)},
				},
			},
			{
				ID: s4ID, ProgramID: progID, SessionName: "Light SBD", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s4ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(3), Reps: intPtr(5), LoadKg: float64Ptr(110), Metadata: mustJSON(map[string]any{"rpe": 6})},
					{ID: uuid.New(), SessionID: s4ID, Order: 2, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(5), LoadKg: float64Ptr(80), Metadata: mustJSON(map[string]any{"rpe": 6})},
					{ID: uuid.New(), SessionID: s4ID, Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(3), Reps: intPtr(5), LoadKg: float64Ptr(140), Metadata: mustJSON(map[string]any{"rpe": 6})},
				},
			},
		},
	}
}

func upperLowerSplit() *domain.Program {
	progID := uuid.New()
	s1ID, s2ID, s3ID, s4ID := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	return &domain.Program{
		ID:     progID,
		Name:   "Upper/Lower Split — Feb 2026",
		Status: domain.ProgramStatusCreated,
		Notes:  stringPtr("Hypertrophy-focused upper/lower split"),
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Upper A", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(6), LoadKg: float64Ptr(90), Metadata: mustJSON(map[string]any{"rpe": 7.5})},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(8), LoadKg: float64Ptr(75)},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(8), LoadKg: float64Ptr(50)},
					{ID: uuid.New(), SessionID: s1ID, Order: 4, ExerciseName: "Barbell Curl", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(30)},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Lower A", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(4), Reps: intPtr(5), LoadKg: float64Ptr(130), Metadata: mustJSON(map[string]any{"rpe": 7.5})},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(8), LoadKg: float64Ptr(100)},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Leg Curl", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(40)},
					{ID: uuid.New(), SessionID: s2ID, Order: 4, ExerciseName: "Calf Raise", Sets: intPtr(4), Reps: intPtr(15), LoadKg: float64Ptr(60)},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Upper B", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Overhead Press", Sets: intPtr(4), Reps: intPtr(5), LoadKg: float64Ptr(55), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Weighted Pull-up", Sets: intPtr(4), Reps: intPtr(6), LoadKg: float64Ptr(20)},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Incline Dumbbell Press", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(30)},
					{ID: uuid.New(), SessionID: s3ID, Order: 4, ExerciseName: "Face Pull", Sets: intPtr(3), Reps: intPtr(15), LoadKg: float64Ptr(20)},
				},
			},
			{
				ID: s4ID, ProgramID: progID, SessionName: "Lower B", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s4ID, Order: 1, ExerciseName: "Conventional Deadlift", Sets: intPtr(4), Reps: intPtr(4), LoadKg: float64Ptr(170), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s4ID, Order: 2, ExerciseName: "Front Squat", Sets: intPtr(3), Reps: intPtr(6), LoadKg: float64Ptr(90)},
					{ID: uuid.New(), SessionID: s4ID, Order: 3, ExerciseName: "Leg Extension", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(50)},
					{ID: uuid.New(), SessionID: s4ID, Order: 4, ExerciseName: "Hip Thrust", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(100)},
				},
			},
		},
	}
}

func blockPeriodization() *domain.Program {
	progID := uuid.New()
	sIDs := [6]uuid.UUID{}
	for i := range sIDs {
		sIDs[i] = uuid.New()
	}

	return &domain.Program{
		ID:     progID,
		Name:   "Block Periodization — Mar 2026",
		Status: domain.ProgramStatusCreated,
		Notes:  stringPtr("Accumulation → Transmutation → Realization"),
		Metadata: mustJSON(map[string]any{
			"block": "accumulation",
			"weeks": 6,
		}),
		Sessions: []domain.ProgramSession{
			{
				ID: sIDs[0], ProgramID: progID, SessionName: "Week 1 — SBD", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[0], Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(4), Reps: intPtr(8), LoadKg: float64Ptr(110), Metadata: mustJSON(map[string]any{"rpe": 7})},
					{ID: uuid.New(), SessionID: sIDs[0], Order: 2, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(8), LoadKg: float64Ptr(80), Metadata: mustJSON(map[string]any{"rpe": 7})},
					{ID: uuid.New(), SessionID: sIDs[0], Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(3), Reps: intPtr(8), LoadKg: float64Ptr(140), Metadata: mustJSON(map[string]any{"rpe": 7})},
				},
			},
			{
				ID: sIDs[1], ProgramID: progID, SessionName: "Week 1 — Accessories", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[1], Order: 1, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(10), LoadKg: float64Ptr(70)},
					{ID: uuid.New(), SessionID: sIDs[1], Order: 2, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(45)},
					{ID: uuid.New(), SessionID: sIDs[1], Order: 3, ExerciseName: "Leg Press", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(180)},
				},
			},
			{
				ID: sIDs[2], ProgramID: progID, SessionName: "Week 2 — SBD", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[2], Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(4), Reps: intPtr(6), LoadKg: float64Ptr(120), Metadata: mustJSON(map[string]any{"rpe": 7.5})},
					{ID: uuid.New(), SessionID: sIDs[2], Order: 2, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(6), LoadKg: float64Ptr(85), Metadata: mustJSON(map[string]any{"rpe": 7.5})},
					{ID: uuid.New(), SessionID: sIDs[2], Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(3), Reps: intPtr(6), LoadKg: float64Ptr(150), Metadata: mustJSON(map[string]any{"rpe": 7.5})},
				},
			},
			{
				ID: sIDs[3], ProgramID: progID, SessionName: "Week 2 — Accessories", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[3], Order: 1, ExerciseName: "Dumbbell Row", Sets: intPtr(4), Reps: intPtr(10), LoadKg: float64Ptr(35)},
					{ID: uuid.New(), SessionID: sIDs[3], Order: 2, ExerciseName: "Dips", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(20), Notes: stringPtr("weighted")},
					{ID: uuid.New(), SessionID: sIDs[3], Order: 3, ExerciseName: "Bulgarian Split Squat", Sets: intPtr(3), Reps: intPtr(10), LoadKg: float64Ptr(40)},
				},
			},
			{
				ID: sIDs[4], ProgramID: progID, SessionName: "Week 3 — SBD", Order: 5,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[4], Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(5), Reps: intPtr(5), LoadKg: float64Ptr(125), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: sIDs[4], Order: 2, ExerciseName: "Bench Press", Sets: intPtr(5), Reps: intPtr(5), LoadKg: float64Ptr(90), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: sIDs[4], Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(4), Reps: intPtr(5), LoadKg: float64Ptr(155), Metadata: mustJSON(map[string]any{"rpe": 8})},
				},
			},
			{
				ID: sIDs[5], ProgramID: progID, SessionName: "Week 3 — Accessories", Order: 6,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[5], Order: 1, ExerciseName: "Pendlay Row", Sets: intPtr(4), Reps: intPtr(8), LoadKg: float64Ptr(75)},
					{ID: uuid.New(), SessionID: sIDs[5], Order: 2, ExerciseName: "Incline Bench Press", Sets: intPtr(3), Reps: intPtr(8), LoadKg: float64Ptr(70)},
					{ID: uuid.New(), SessionID: sIDs[5], Order: 3, ExerciseName: "Leg Curl", Sets: intPtr(3), Reps: intPtr(12), LoadKg: float64Ptr(45)},
				},
			},
		},
	}
}

func competitionPrep() *domain.Program {
	progID := uuid.New()
	s1ID, s2ID, s3ID := uuid.New(), uuid.New(), uuid.New()

	return &domain.Program{
		ID:     progID,
		Name:   "Competition Prep — Apr 2026",
		Status: domain.ProgramStatusCreated,
		Notes:  stringPtr("Final peaking block before April competition"),
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Opener Rehearsal", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(3), Reps: intPtr(1), LoadKg: float64Ptr(150), Notes: stringPtr("opener attempt"), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(1), LoadKg: float64Ptr(107.5), Notes: stringPtr("opener attempt"), Metadata: mustJSON(map[string]any{"rpe": 8})},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(3), Reps: intPtr(1), LoadKg: float64Ptr(195), Notes: stringPtr("opener attempt"), Metadata: mustJSON(map[string]any{"rpe": 8})},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Speed Work", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(6), Reps: intPtr(2), LoadKg: float64Ptr(110), Notes: stringPtr("60% 1RM, focus on speed"), Metadata: mustJSON(map[string]any{"rpe": 6})},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Bench Press", Sets: intPtr(6), Reps: intPtr(2), LoadKg: float64Ptr(75), Notes: stringPtr("60% 1RM, focus on speed"), Metadata: mustJSON(map[string]any{"rpe": 6})},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(5), Reps: intPtr(2), LoadKg: float64Ptr(135), Notes: stringPtr("60% 1RM, focus on speed"), Metadata: mustJSON(map[string]any{"rpe": 6})},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Light Recovery", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Low Bar Squat", Sets: intPtr(3), Reps: intPtr(3), LoadKg: float64Ptr(90), Metadata: mustJSON(map[string]any{"rpe": 5})},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(3), LoadKg: float64Ptr(65), Metadata: mustJSON(map[string]any{"rpe": 5})},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Conventional Deadlift", Sets: intPtr(2), Reps: intPtr(3), LoadKg: float64Ptr(110), Metadata: mustJSON(map[string]any{"rpe": 5})},
				},
			},
		},
	}
}

// logForSession creates a Log from a ProgramSession, copying exercise data.
func logForSession(programID uuid.UUID, session domain.ProgramSession, performedAt time.Time) *domain.Log {
	logID := uuid.New()
	entries := make([]domain.LogEntry, len(session.Entries))
	for i, e := range session.Entries {
		entries[i] = domain.LogEntry{
			ID:           uuid.New(),
			LogID:        logID,
			Order:        e.Order,
			ExerciseName: e.ExerciseName,
			Sets:         e.Sets,
			Reps:         e.Reps,
			LoadKg:       e.LoadKg,
			Notes:        e.Notes,
			Metadata:     e.Metadata,
		}
	}

	sessionName := session.SessionName
	return &domain.Log{
		ID:          logID,
		ProgramID:   &programID,
		SessionName: &sessionName,
		PerformedAt: performedAt,
		Entries:     entries,
	}
}
