// Package seed provides local development seed data.
// It is intended for use with in-memory storage and local postgres only.
//
// NOTE: This file is a stub pending Task 14 (Final Integration & Cleanup).
// The seed data needs to be updated for the new domain model (no Status/Metadata on Program).
package seed

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

func stringPtr(v string) *string { return &v }

func fields(sets, reps int, loadKg float64) map[string]interface{} {
	return map[string]interface{}{"sets": sets, "reps": reps, "load_kg": loadKg}
}

func fieldsWithRPE(sets, reps int, loadKg, rpe float64) map[string]interface{} {
	return map[string]interface{}{"sets": sets, "reps": reps, "load_kg": loadKg, "rpe": rpe}
}

// Run inserts sample powerlifting data into the provided repositories.
func Run(ctx context.Context, programRepo repository.ProgramRepository, logRepo repository.LogRepository) error {
	// ── Program 1: SBD Peaking ──────────────────────────────────
	prog1 := sbdPeaking()
	if err := programRepo.Create(ctx, prog1); err != nil {
		return fmt.Errorf("create SBD Peaking program: %w", err)
	}
	slog.Info("seeded program", "name", prog1.Name)

	// Logs for all 4 sessions
	for i, s := range prog1.Sessions {
		log := logForSession(prog1.ID, s, time.Date(2026, 1, 6+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog1.Name, s.SessionName, err)
		}
	}

	// ── Program 2: Upper/Lower Split ────────────────────────────
	prog2 := upperLowerSplit()
	if err := programRepo.Create(ctx, prog2); err != nil {
		return fmt.Errorf("create Upper/Lower Split program: %w", err)
	}
	slog.Info("seeded program", "name", prog2.Name)

	for i, s := range prog2.Sessions {
		log := logForSession(prog2.ID, s, time.Date(2026, 2, 2+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog2.Name, s.SessionName, err)
		}
	}

	// ── Program 3: Block Periodization (3/6 logged) ────────────────
	prog3 := blockPeriodization()
	if err := programRepo.Create(ctx, prog3); err != nil {
		return fmt.Errorf("create Block Periodization program: %w", err)
	}
	slog.Info("seeded program", "name", prog3.Name)

	for i := 0; i < 3; i++ {
		s := prog3.Sessions[i]
		log := logForSession(prog3.ID, s, time.Date(2026, 3, 3+i*2, 18, 0, 0, 0, time.UTC))
		if err := logRepo.Create(ctx, log); err != nil {
			return fmt.Errorf("create log for %s/%s: %w", prog3.Name, s.SessionName, err)
		}
	}

	// ── Program 4: Competition Prep (no logs) ──────────────────────
	prog4 := competitionPrep()
	if err := programRepo.Create(ctx, prog4); err != nil {
		return fmt.Errorf("create Competition Prep program: %w", err)
	}
	slog.Info("seeded program", "name", prog4.Name)

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
		Notes:  stringPtr("4-session peaking block focusing on competition lifts"),
		Status: domain.ProgramStatusPublished,
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Heavy Squat", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(5, 3, 140, 8)},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Pause Squat", Fields: fieldsWithRPE(3, 3, 120, 7)},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Leg Press", Fields: fields(3, 10, 200)},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Heavy Bench", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Bench Press", Fields: fieldsWithRPE(5, 3, 100, 8)},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Close Grip Bench Press", Fields: fieldsWithRPE(3, 5, 85, 7)},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Dumbbell Fly", Fields: fields(3, 12, 16)},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Heavy Deadlift", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(5, 2, 180, 8.5)},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Deficit Deadlift", Fields: fieldsWithRPE(3, 4, 150, 7)},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Barbell Row", Fields: fields(4, 8, 80)},
				},
			},
			{
				ID: s4ID, ProgramID: progID, SessionName: "Light SBD", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s4ID, Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(3, 5, 110, 6)},
					{ID: uuid.New(), SessionID: s4ID, Order: 2, ExerciseName: "Bench Press", Fields: fieldsWithRPE(3, 5, 80, 6)},
					{ID: uuid.New(), SessionID: s4ID, Order: 3, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(3, 5, 140, 6)},
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
		Notes:  stringPtr("Hypertrophy-focused upper/lower split"),
		Status: domain.ProgramStatusPublished,
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Upper A", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Bench Press", Fields: fieldsWithRPE(4, 6, 90, 7.5)},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Barbell Row", Fields: fields(4, 8, 75)},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Overhead Press", Fields: fields(3, 8, 50)},
					{ID: uuid.New(), SessionID: s1ID, Order: 4, ExerciseName: "Barbell Curl", Fields: fields(3, 12, 30)},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Lower A", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(4, 5, 130, 7.5)},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Romanian Deadlift", Fields: fields(3, 8, 100)},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Leg Curl", Fields: fields(3, 12, 40)},
					{ID: uuid.New(), SessionID: s2ID, Order: 4, ExerciseName: "Calf Raise", Fields: fields(4, 15, 60)},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Upper B", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Overhead Press", Fields: fieldsWithRPE(4, 5, 55, 8)},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Weighted Pull-up", Fields: fields(4, 6, 20)},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Incline Dumbbell Press", Fields: fields(3, 10, 30)},
					{ID: uuid.New(), SessionID: s3ID, Order: 4, ExerciseName: "Face Pull", Fields: fields(3, 15, 20)},
				},
			},
			{
				ID: s4ID, ProgramID: progID, SessionName: "Lower B", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s4ID, Order: 1, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(4, 4, 170, 8)},
					{ID: uuid.New(), SessionID: s4ID, Order: 2, ExerciseName: "Front Squat", Fields: fields(3, 6, 90)},
					{ID: uuid.New(), SessionID: s4ID, Order: 3, ExerciseName: "Leg Extension", Fields: fields(3, 12, 50)},
					{ID: uuid.New(), SessionID: s4ID, Order: 4, ExerciseName: "Hip Thrust", Fields: fields(3, 10, 100)},
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
		Notes:  stringPtr("Accumulation → Transmutation → Realization"),
		Status: domain.ProgramStatusPublished,
		Sessions: []domain.ProgramSession{
			{
				ID: sIDs[0], ProgramID: progID, SessionName: "Week 1 — SBD", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[0], Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(4, 8, 110, 7)},
					{ID: uuid.New(), SessionID: sIDs[0], Order: 2, ExerciseName: "Bench Press", Fields: fieldsWithRPE(4, 8, 80, 7)},
					{ID: uuid.New(), SessionID: sIDs[0], Order: 3, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(3, 8, 140, 7)},
				},
			},
			{
				ID: sIDs[1], ProgramID: progID, SessionName: "Week 1 — Accessories", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[1], Order: 1, ExerciseName: "Barbell Row", Fields: fields(4, 10, 70)},
					{ID: uuid.New(), SessionID: sIDs[1], Order: 2, ExerciseName: "Overhead Press", Fields: fields(3, 10, 45)},
					{ID: uuid.New(), SessionID: sIDs[1], Order: 3, ExerciseName: "Leg Press", Fields: fields(3, 12, 180)},
				},
			},
			{
				ID: sIDs[2], ProgramID: progID, SessionName: "Week 2 — SBD", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[2], Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(4, 6, 120, 7.5)},
					{ID: uuid.New(), SessionID: sIDs[2], Order: 2, ExerciseName: "Bench Press", Fields: fieldsWithRPE(4, 6, 85, 7.5)},
					{ID: uuid.New(), SessionID: sIDs[2], Order: 3, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(3, 6, 150, 7.5)},
				},
			},
			{
				ID: sIDs[3], ProgramID: progID, SessionName: "Week 2 — Accessories", Order: 4,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[3], Order: 1, ExerciseName: "Dumbbell Row", Fields: fields(4, 10, 35)},
					{ID: uuid.New(), SessionID: sIDs[3], Order: 2, ExerciseName: "Dips", Fields: fields(3, 10, 20), Notes: stringPtr("weighted")},
					{ID: uuid.New(), SessionID: sIDs[3], Order: 3, ExerciseName: "Bulgarian Split Squat", Fields: fields(3, 10, 40)},
				},
			},
			{
				ID: sIDs[4], ProgramID: progID, SessionName: "Week 3 — SBD", Order: 5,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[4], Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(5, 5, 125, 8)},
					{ID: uuid.New(), SessionID: sIDs[4], Order: 2, ExerciseName: "Bench Press", Fields: fieldsWithRPE(5, 5, 90, 8)},
					{ID: uuid.New(), SessionID: sIDs[4], Order: 3, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(4, 5, 155, 8)},
				},
			},
			{
				ID: sIDs[5], ProgramID: progID, SessionName: "Week 3 — Accessories", Order: 6,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: sIDs[5], Order: 1, ExerciseName: "Pendlay Row", Fields: fields(4, 8, 75)},
					{ID: uuid.New(), SessionID: sIDs[5], Order: 2, ExerciseName: "Incline Bench Press", Fields: fields(3, 8, 70)},
					{ID: uuid.New(), SessionID: sIDs[5], Order: 3, ExerciseName: "Leg Curl", Fields: fields(3, 12, 45)},
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
		Notes:  stringPtr("Final peaking block before April competition"),
		Status: domain.ProgramStatusPublished,
		Sessions: []domain.ProgramSession{
			{
				ID: s1ID, ProgramID: progID, SessionName: "Opener Rehearsal", Order: 1,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s1ID, Order: 1, ExerciseName: "Low Bar Squat", Notes: stringPtr("opener attempt"), Fields: fieldsWithRPE(3, 1, 150, 8)},
					{ID: uuid.New(), SessionID: s1ID, Order: 2, ExerciseName: "Bench Press", Notes: stringPtr("opener attempt"), Fields: fieldsWithRPE(3, 1, 107.5, 8)},
					{ID: uuid.New(), SessionID: s1ID, Order: 3, ExerciseName: "Conventional Deadlift", Notes: stringPtr("opener attempt"), Fields: fieldsWithRPE(3, 1, 195, 8)},
				},
			},
			{
				ID: s2ID, ProgramID: progID, SessionName: "Speed Work", Order: 2,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s2ID, Order: 1, ExerciseName: "Low Bar Squat", Notes: stringPtr("60% 1RM, focus on speed"), Fields: fieldsWithRPE(6, 2, 110, 6)},
					{ID: uuid.New(), SessionID: s2ID, Order: 2, ExerciseName: "Bench Press", Notes: stringPtr("60% 1RM, focus on speed"), Fields: fieldsWithRPE(6, 2, 75, 6)},
					{ID: uuid.New(), SessionID: s2ID, Order: 3, ExerciseName: "Conventional Deadlift", Notes: stringPtr("60% 1RM, focus on speed"), Fields: fieldsWithRPE(5, 2, 135, 6)},
				},
			},
			{
				ID: s3ID, ProgramID: progID, SessionName: "Light Recovery", Order: 3,
				Entries: []domain.ProgramSessionEntry{
					{ID: uuid.New(), SessionID: s3ID, Order: 1, ExerciseName: "Low Bar Squat", Fields: fieldsWithRPE(3, 3, 90, 5)},
					{ID: uuid.New(), SessionID: s3ID, Order: 2, ExerciseName: "Bench Press", Fields: fieldsWithRPE(3, 3, 65, 5)},
					{ID: uuid.New(), SessionID: s3ID, Order: 3, ExerciseName: "Conventional Deadlift", Fields: fieldsWithRPE(2, 3, 110, 5)},
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
			Fields:       e.Fields,
			Notes:        e.Notes,
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
