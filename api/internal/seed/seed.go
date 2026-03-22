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
)

// Run inserts sample powerlifting data into the provided repositories.
//
// Program lifecycle simulated:
//
//	Completed 1: SBD Peaking — Jan 2026       (Template 1, light weights, beginner)
//	Completed 2: Upper/Lower — Feb 2026        (Template 2, medium weights)
//	Completed 3: Block Periodization — Mar 2026 (Template 4, first block attempt)
//	Active:      SBD Peaking — Mar 2026        (Template 1, current, 3/6 sessions done)
//	Planned:     Competition Peak — Apr 2026   (Template 5, not started)
func Run(ctx context.Context, programTemplateRepo repository.ProgramTemplateRepository, programRepo repository.ProgramRepository, logRepo repository.LogRepository) error {
	// ── Create 5 templates ────────────────────────────────────────────────────

	sbdPeakTmpl := sbdPeakingTemplate()
	if err := programTemplateRepo.Create(ctx, sbdPeakTmpl); err != nil {
		return fmt.Errorf("create SBD peaking template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", sbdPeakTmpl.Name)

	ulTmpl := upperLower4WeekTemplate()
	if err := programTemplateRepo.Create(ctx, ulTmpl); err != nil {
		return fmt.Errorf("create upper/lower template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", ulTmpl.Name)

	pplTmpl := pushPullLegs3WeekTemplate()
	if err := programTemplateRepo.Create(ctx, pplTmpl); err != nil {
		return fmt.Errorf("create PPL template: %w", err)
	}
	if err := programTemplateRepo.Archive(ctx, pplTmpl.ID); err != nil {
		return fmt.Errorf("archive PPL template: %w", err)
	}
	slog.Info("[seed] Created and archived program template", "name", pplTmpl.Name)

	blockTmpl := sbdBlockPeriodizationTemplate()
	if err := programTemplateRepo.Create(ctx, blockTmpl); err != nil {
		return fmt.Errorf("create block periodization template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", blockTmpl.Name)

	compPeakTmpl := sbdCompetitionPeakTemplate()
	if err := programTemplateRepo.Create(ctx, compPeakTmpl); err != nil {
		return fmt.Errorf("create competition peak template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", compPeakTmpl.Name)

	// ── Completed 1: SBD Peaking — January 2026 (Template 1, beginner weights) ──

	comp1 := domain.GenerateProgramFromTemplate(sbdPeakTmpl, &domain.GenerateProgramInput{
		Name:           "SBD Peaking — January 2026",
		TargetWeights:  map[string]float64{"Squat": 110.0, "Bench": 82.5, "Deadlift": 140.0},
		LoadIncrements: map[string]float64{"Squat": 2.5, "Bench": 2.5, "Deadlift": 2.5},
	})
	if err := programRepo.Create(ctx, comp1); err != nil {
		return fmt.Errorf("create comp1 program: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, comp1.ID, domain.ProgramStatusCompleted); err != nil {
		return fmt.Errorf("complete comp1 program: %w", err)
	}
	slog.Info("[seed] Created completed program", "name", comp1.Name)

	// ── Completed 2: Upper/Lower — February 2026 (Template 2) ──

	comp2 := domain.GenerateProgramFromTemplate(ulTmpl, &domain.GenerateProgramInput{
		Name: "Upper/Lower Strength — February 2026",
		TargetWeights: map[string]float64{
			"Squat": 115.0, "Deadlift": 145.0, "Romanian Deadlift": 90.0,
			"Bench Press": 87.5, "Overhead Press": 60.0, "Barbell Row": 75.0,
		},
		LoadIncrements: map[string]float64{
			"Squat": 2.5, "Deadlift": 2.5, "Romanian Deadlift": 2.5,
			"Bench Press": 2.5, "Overhead Press": 2.5, "Barbell Row": 2.5,
		},
	})
	if err := programRepo.Create(ctx, comp2); err != nil {
		return fmt.Errorf("create comp2 program: %w", err)
	}
	if err := programRepo.UpdateStatus(ctx, comp2.ID, domain.ProgramStatusCompleted); err != nil {
		return fmt.Errorf("complete comp2 program: %w", err)
	}
	slog.Info("[seed] Created completed program", "name", comp2.Name)

	// ── Completed 3: SBD Block Periodization — Early March 2026 (Template 4) ──

	comp3 := domain.GenerateProgramFromTemplate(blockTmpl, &domain.GenerateProgramInput{
		Name:           "SBD Block Periodization — Q1 2026",
		TargetWeights:  map[string]float64{"Squat": 130.0, "Bench": 97.5, "Deadlift": 162.5},
		LoadIncrements: map[string]float64{"Squat": 2.5, "Bench": 2.5, "Deadlift": 2.5},
	})
	if err := programRepo.Create(ctx, comp3); err != nil {
		return fmt.Errorf("create comp3 program: %w", err)
	}
	comp3Logs := blockPeriodizationLogs(comp3)
	for _, l := range comp3Logs {
		lCopy := l
		if err := logRepo.Create(ctx, &lCopy); err != nil {
			return fmt.Errorf("create block log %s: %w", l.PerformedAt.Format("2006-01-02"), err)
		}
		slog.Info("[seed] Created log", "date", l.PerformedAt.Format("2006-01-02"))
	}
	if err := programRepo.UpdateStatus(ctx, comp3.ID, domain.ProgramStatusCompleted); err != nil {
		return fmt.Errorf("complete comp3 program: %w", err)
	}
	slog.Info("[seed] Created completed program", "name", comp3.Name)

	// ── Active: SBD Peaking — March 2026 (Template 1, current cycle) ──

	activeProgram := domain.GenerateProgramFromTemplate(sbdPeakTmpl, &domain.GenerateProgramInput{
		Name:           "SBD Peaking — March 2026",
		TargetWeights:  map[string]float64{"Squat": 145.0, "Bench": 110.0, "Deadlift": 180.0},
		LoadIncrements: map[string]float64{"Squat": 2.5, "Bench": 2.5, "Deadlift": 2.5},
	})
	if err := programRepo.Create(ctx, activeProgram); err != nil {
		return fmt.Errorf("create active program: %w", err)
	}
	slog.Info("[seed] Created active program", "name", activeProgram.Name)

	for _, l := range activeProgramLogs(activeProgram) {
		lCopy := l
		if err := logRepo.Create(ctx, &lCopy); err != nil {
			return fmt.Errorf("create active log %s: %w", l.PerformedAt.Format("2006-01-02"), err)
		}
		slog.Info("[seed] Created log", "date", l.PerformedAt.Format("2006-01-02"))
	}

	// ── Planned: Competition Peak — April 2026 (Template 5) ──

	plannedProgram := domain.GenerateProgramFromTemplate(compPeakTmpl, &domain.GenerateProgramInput{
		Name:           "Competition Peak — April 2026",
		TargetWeights:  map[string]float64{"Squat": 150.0, "Bench": 115.0, "Deadlift": 190.0},
		LoadIncrements: map[string]float64{"Squat": 2.5, "Bench": 2.5, "Deadlift": 2.5},
	})
	plannedProgram.Status = domain.ProgramStatusPlanned
	plannedProgram.Notes = strPtr("Target: SQ 150 / BP 115 / DL 190 — competing at April meet")
	if err := programRepo.Create(ctx, plannedProgram); err != nil {
		return fmt.Errorf("create planned program: %w", err)
	}
	slog.Info("[seed] Created planned program", "name", plannedProgram.Name)

	return nil
}

// ── Template definitions ─────────────────────────────────────────────────────

// sbdPeakingTemplate is a 3-week SBD peaking cycle with week-only session structure.
// Sessions: Week 1 Day 1 / Day 2, Week 2 Day 1 / Day 2, Week 3 Day 1 / Day 2
func sbdPeakingTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "SBD 3-Week Peaking",
		Description: strPtr("3-week Squat / Bench / Deadlift peaking cycle (2 days/week). Intensity rises each week toward a peak single."),
		Entries: []domain.ProgramTemplateEntry{
			// ── Week 1 Day 1: Heavy ──
			{Order: 1, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 2, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			{Order: 3, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
			{Order: 4, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 5, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			{Order: 6, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
			{Order: 7, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 8, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			// ── Week 1 Day 2: Volume ──
			{Order: 9, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			{Order: 10, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			{Order: 11, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			// ── Week 2 Day 1: Heavy ──
			{Order: 12, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 13, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			{Order: 14, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
			{Order: 15, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 16, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			{Order: 17, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
			{Order: 18, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 19, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			// ── Week 2 Day 2: Volume ──
			{Order: 20, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			{Order: 21, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			{Order: 22, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			// ── Week 3 Day 1: Peak ──
			{Order: 23, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 24, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			{Order: 25, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
			{Order: 26, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 27, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			{Order: 28, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
			{Order: 29, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 30, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			// ── Week 3 Day 2: Deload ──
			{Order: 31, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
			{Order: 32, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
			{Order: 33, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
		},
	}
}

// upperLower4WeekTemplate is a 4-week upper/lower split with week-only session structure.
// Sessions: Week N Upper / Week N Lower (8 sessions total).
// Volume decreases and intensity increases week over week.
func upperLower4WeekTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "Upper/Lower 4-Week Strength",
		Description: strPtr("4-week upper/lower split (2 days/week). Volume decreases and intensity increases each week."),
		Entries: []domain.ProgramTemplateEntry{
			// ── Week 1 Upper (Volume) ──
			{Order: 1, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Upper")},
			{Order: 2, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Upper")},
			{Order: 3, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Upper")},
			// ── Week 1 Lower (Volume) ──
			{Order: 4, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Lower")},
			{Order: 5, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Lower")},
			{Order: 6, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Lower")},
			// ── Week 2 Upper ──
			{Order: 7, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Upper")},
			{Order: 8, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Upper")},
			{Order: 9, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Upper")},
			// ── Week 2 Lower ──
			{Order: 10, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Lower")},
			{Order: 11, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Lower")},
			{Order: 12, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Lower")},
			// ── Week 3 Upper ──
			{Order: 13, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Upper")},
			{Order: 14, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Upper")},
			{Order: 15, ExerciseName: "Barbell Row", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Upper")},
			// ── Week 3 Lower ──
			{Order: 16, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Lower")},
			{Order: 17, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Lower")},
			{Order: 18, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Lower")},
			// ── Week 4 Upper (Peak) ──
			{Order: 19, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Week 4 Upper")},
			{Order: 20, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Week 4 Upper")},
			{Order: 21, ExerciseName: "Barbell Row", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Week 4 Upper")},
			// ── Week 4 Lower (Peak) ──
			{Order: 22, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Week 4 Lower")},
			{Order: 23, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Week 4 Lower")},
			{Order: 24, ExerciseName: "Romanian Deadlift", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(8), Metadata: sessionMeta("Week 4 Lower")},
		},
	}
}

// pushPullLegs3WeekTemplate is a 3-week PPL split with week-only session structure.
// This template is archived — kept for historical reference.
// Sessions: Week N Push / Pull / Legs (9 sessions total).
func pushPullLegs3WeekTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "Push/Pull/Legs 3-Week",
		Description: strPtr("3-week push/pull/legs split (3 days/week). Archived — superseded by Upper/Lower program."),
		Entries: []domain.ProgramTemplateEntry{
			// ── Week 1 Push ──
			{Order: 1, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Push")},
			{Order: 2, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Push")},
			{Order: 3, ExerciseName: "Tricep Pushdown", Sets: intPtr(3), Reps: intPtr(12), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Push")},
			// ── Week 1 Pull ──
			{Order: 4, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Pull")},
			{Order: 5, ExerciseName: "Pull-up", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Pull")},
			{Order: 6, ExerciseName: "Bicep Curl", Sets: intPtr(3), Reps: intPtr(12), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Pull")},
			// ── Week 1 Legs ──
			{Order: 7, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Legs")},
			{Order: 8, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Legs")},
			{Order: 9, ExerciseName: "Leg Press", Sets: intPtr(3), Reps: intPtr(12), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Legs")},
			// ── Week 2 Push ──
			{Order: 10, ExerciseName: "Bench Press", Sets: intPtr(4), Reps: intPtr(7), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Push")},
			{Order: 11, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Push")},
			{Order: 12, ExerciseName: "Tricep Pushdown", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Push")},
			// ── Week 2 Pull ──
			{Order: 13, ExerciseName: "Barbell Row", Sets: intPtr(4), Reps: intPtr(7), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Pull")},
			{Order: 14, ExerciseName: "Pull-up", Sets: intPtr(3), Reps: intPtr(7), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Pull")},
			{Order: 15, ExerciseName: "Bicep Curl", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Pull")},
			// ── Week 2 Legs ──
			{Order: 16, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Legs")},
			{Order: 17, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Legs")},
			{Order: 18, ExerciseName: "Leg Press", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Legs")},
			// ── Week 3 Push ──
			{Order: 19, ExerciseName: "Bench Press", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(9), Metadata: sessionMeta("Week 3 Push")},
			{Order: 20, ExerciseName: "Overhead Press", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Push")},
			{Order: 21, ExerciseName: "Tricep Pushdown", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Push")},
			// ── Week 3 Pull ──
			{Order: 22, ExerciseName: "Barbell Row", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(9), Metadata: sessionMeta("Week 3 Pull")},
			{Order: 23, ExerciseName: "Pull-up", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Pull")},
			{Order: 24, ExerciseName: "Bicep Curl", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Pull")},
			// ── Week 3 Legs ──
			{Order: 25, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(9), Metadata: sessionMeta("Week 3 Legs")},
			{Order: 26, ExerciseName: "Romanian Deadlift", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Legs")},
			{Order: 27, ExerciseName: "Leg Press", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 3 Legs")},
		},
	}
}

// sbdBlockPeriodizationTemplate is an 8-week SBD program split into two 4-week blocks.
// Block 1 (Accumulation): higher volume, RPE 7–8.
// Block 2 (Intensification): lower volume, RPE 8–9.
// Sessions use block-prefixed names, e.g. "Accumulation W1 D1".
func sbdBlockPeriodizationTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "SBD 8-Week Block Periodization",
		Description: strPtr("8-week SBD program with 4-week blocks: Accumulation (volume) → Intensification (strength). 2 days/week per block."),
		Entries: []domain.ProgramTemplateEntry{
			// ── Block 1: Accumulation ──────────────────────────────
			// Accumulation W1 D1
			{Order: 1, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D1")},
			{Order: 2, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D1")},
			{Order: 3, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D1")},
			// Accumulation W1 D2
			{Order: 4, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D2")},
			{Order: 5, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D2")},
			{Order: 6, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W1 D2")},
			// Accumulation W2 D1
			{Order: 7, ExerciseName: "Squat", Sets: intPtr(5), Reps: intPtr(7), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D1")},
			{Order: 8, ExerciseName: "Bench", Sets: intPtr(5), Reps: intPtr(7), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D1")},
			{Order: 9, ExerciseName: "Deadlift", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D1")},
			// Accumulation W2 D2
			{Order: 10, ExerciseName: "Squat", Sets: intPtr(5), Reps: intPtr(7), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D2")},
			{Order: 11, ExerciseName: "Bench", Sets: intPtr(5), Reps: intPtr(7), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D2")},
			{Order: 12, ExerciseName: "Deadlift", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionMeta("Accumulation W2 D2")},
			// Accumulation W3 D1
			{Order: 13, ExerciseName: "Squat", Sets: intPtr(5), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D1")},
			{Order: 14, ExerciseName: "Bench", Sets: intPtr(5), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D1")},
			{Order: 15, ExerciseName: "Deadlift", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D1")},
			// Accumulation W3 D2
			{Order: 16, ExerciseName: "Squat", Sets: intPtr(5), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D2")},
			{Order: 17, ExerciseName: "Bench", Sets: intPtr(5), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D2")},
			{Order: 18, ExerciseName: "Deadlift", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Accumulation W3 D2")},
			// Accumulation W4 D1 (deload)
			{Order: 19, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D1")},
			{Order: 20, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D1")},
			{Order: 21, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D1")},
			// Accumulation W4 D2 (deload)
			{Order: 22, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D2")},
			{Order: 23, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D2")},
			{Order: 24, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionMeta("Accumulation W4 D2")},
			// ── Block 2: Intensification ──────────────────────────
			// Intensification W1 D1
			{Order: 25, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D1")},
			{Order: 26, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D1")},
			{Order: 27, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D1")},
			// Intensification W1 D2
			{Order: 28, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D2")},
			{Order: 29, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D2")},
			{Order: 30, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionMeta("Intensification W1 D2")},
			// Intensification W2 D1
			{Order: 31, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D1")},
			{Order: 32, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D1")},
			{Order: 33, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D1")},
			// Intensification W2 D2
			{Order: 34, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D2")},
			{Order: 35, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D2")},
			{Order: 36, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W2 D2")},
			// Intensification W3 D1
			{Order: 37, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D1")},
			{Order: 38, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D1")},
			{Order: 39, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D1")},
			// Intensification W3 D2
			{Order: 40, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D2")},
			{Order: 41, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D2")},
			{Order: 42, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Intensification W3 D2")},
			// Intensification W4 D1 (peak singles)
			{Order: 43, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Intensification W4 D1", "top")},
			{Order: 44, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Intensification W4 D1", "top")},
			{Order: 45, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Intensification W4 D1", "top")},
			// Intensification W4 D2 (deload)
			{Order: 46, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(6), Metadata: sessionMeta("Intensification W4 D2")},
			{Order: 47, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(6), Metadata: sessionMeta("Intensification W4 D2")},
			{Order: 48, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(6), Metadata: sessionMeta("Intensification W4 D2")},
		},
	}
}

// sbdCompetitionPeakTemplate is a 4-week competition preparation block.
// All sessions belong to a single "Peak Block", grouping 4 weeks into one named block.
// Session names: "Peak Block W1 D1" … "Peak Block W4 D2"
func sbdCompetitionPeakTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "SBD 4-Week Competition Peak",
		Description: strPtr("4-week competition peak block (2 days/week). Designed to peak strength for meet day. All 4 weeks form a single Peak Block."),
		Entries: []domain.ProgramTemplateEntry{
			// ── Peak Block W1 ──
			{Order: 1, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D1")},
			{Order: 2, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D1")},
			{Order: 3, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D1")},
			{Order: 4, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D2")},
			{Order: 5, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D2")},
			{Order: 6, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W1 D2")},
			// ── Peak Block W2 ──
			{Order: 7, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W2 D1")},
			{Order: 8, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W2 D1")},
			{Order: 9, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionMeta("Peak Block W2 D1")},
			{Order: 10, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W2 D2")},
			{Order: 11, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W2 D2")},
			{Order: 12, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W2 D2")},
			// ── Peak Block W3 ──
			{Order: 13, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D1")},
			{Order: 14, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D1")},
			{Order: 15, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D1")},
			{Order: 16, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D2")},
			{Order: 17, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D2")},
			{Order: 18, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionMeta("Peak Block W3 D2")},
			// ── Peak Block W4 (meet week) ──
			{Order: 19, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Peak Block W4 D1", "opener")},
			{Order: 20, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Peak Block W4 D1", "opener")},
			{Order: 21, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Peak Block W4 D1", "opener")},
			{Order: 22, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(6), Metadata: sessionLabelMeta("Peak Block W4 D2", "comp")},
			{Order: 23, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(6), Metadata: sessionLabelMeta("Peak Block W4 D2", "comp")},
			{Order: 24, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(6), Metadata: sessionLabelMeta("Peak Block W4 D2", "comp")},
		},
	}
}

// ── Log definitions ───────────────────────────────────────────────────────────

// blockPeriodizationLogs returns 2 representative logs from the Accumulation block.
// The rest of the sessions are implied via UpdateStatus(completed).
func blockPeriodizationLogs(program *domain.Program) []domain.Log {
	programID := program.ID
	s1 := "Accumulation W1 D1"
	s2 := "Accumulation W1 D2"

	return []domain.Log{
		{
			ProgramID:   &programID,
			SessionName: &s1,
			PerformedAt: time.Date(2026, 1, 20, 18, 0, 0, 0, time.Local),
			Notes:       strPtr("Accumulation W1 D1 — volume felt manageable"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(105), RPE: intPtr(7)},
				{Order: 2, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(77.5), RPE: intPtr(7)},
				{Order: 3, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(6), LoadKg: f64Ptr(130), RPE: intPtr(7)},
			},
		},
		{
			ProgramID:   &programID,
			SessionName: &s2,
			PerformedAt: time.Date(2026, 1, 23, 18, 0, 0, 0, time.Local),
			Notes:       strPtr("Accumulation W1 D2 — squats felt heavier than expected"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(105), RPE: intPtr(8)},
				{Order: 2, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(77.5), RPE: intPtr(7)},
				{Order: 3, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(6), LoadKg: f64Ptr(130), RPE: intPtr(7)},
			},
		},
	}
}

// activeProgramLogs returns logs for the first 3 sessions of the current SBD peaking program.
func activeProgramLogs(program *domain.Program) []domain.Log {
	programID := program.ID
	w1d1 := "Week 1 Day 1"
	w1d2 := "Week 1 Day 2"
	w2d1 := "Week 2 Day 1"

	return []domain.Log{
		// Week 1 Day 1: Heavy
		{
			ProgramID:   &programID,
			SessionName: &w1d1,
			PerformedAt: time.Date(2026, 2, 26, 18, 0, 0, 0, time.Local),
			StartedAt:   timePtr(time.Date(2026, 2, 26, 18, 0, 0, 0, time.Local)),
			FinishedAt:  timePtr(time.Date(2026, 2, 26, 19, 22, 0, 0, time.Local)),
			Notes:       strPtr("Week 1 heavy — felt solid"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(122.5), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 2, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(105), RPE: intPtr(8), Metadata: labelMeta("main")},
				{Order: 3, ExerciseName: "Squat", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(95), RPE: intPtr(7), Metadata: labelMeta("backoff")},
				{Order: 4, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(92.5), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 5, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(80), RPE: intPtr(8), Metadata: labelMeta("main")},
				{Order: 6, ExerciseName: "Bench", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(72.5), RPE: intPtr(7), Metadata: labelMeta("backoff")},
				{Order: 7, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(150), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 8, ExerciseName: "Deadlift", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(130), RPE: intPtr(8), Metadata: labelMeta("main")},
			},
		},
		// Week 1 Day 2: Volume
		{
			ProgramID:   &programID,
			SessionName: &w1d2,
			PerformedAt: time.Date(2026, 2, 28, 18, 0, 0, 0, time.Local),
			StartedAt:   timePtr(time.Date(2026, 2, 28, 18, 0, 0, 0, time.Local)),
			FinishedAt:  timePtr(time.Date(2026, 2, 28, 18, 48, 0, 0, time.Local)),
			Notes:       strPtr("Week 1 volume — kept technique clean"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(97.5), RPE: intPtr(7), Metadata: labelMeta("main")},
				{Order: 2, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(72.5), RPE: intPtr(7), Metadata: labelMeta("main")},
				{Order: 3, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(125), RPE: intPtr(7), Metadata: labelMeta("main")},
			},
		},
		// Week 2 Day 1: Heavy
		{
			ProgramID:   &programID,
			SessionName: &w2d1,
			PerformedAt: time.Date(2026, 3, 5, 18, 0, 0, 0, time.Local),
			StartedAt:   timePtr(time.Date(2026, 3, 5, 18, 0, 0, 0, time.Local)),
			FinishedAt:  timePtr(time.Date(2026, 3, 5, 19, 35, 0, 0, time.Local)),
			Notes:       strPtr("Week 2 heavy — bench feeling great"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "Squat", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(125), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 2, ExerciseName: "Squat", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(107.5), RPE: intPtr(8), Metadata: labelMeta("main")},
				{Order: 3, ExerciseName: "Squat", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(97.5), RPE: intPtr(7), Metadata: labelMeta("backoff")},
				{Order: 4, ExerciseName: "Bench", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(95), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 5, ExerciseName: "Bench", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(82.5), RPE: intPtr(8), Metadata: labelMeta("main")},
				{Order: 6, ExerciseName: "Bench", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(75), RPE: intPtr(7), Metadata: labelMeta("backoff")},
				{Order: 7, ExerciseName: "Deadlift", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(155), RPE: intPtr(9), Metadata: labelMeta("top")},
				{Order: 8, ExerciseName: "Deadlift", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(135), RPE: intPtr(8), Metadata: labelMeta("main")},
			},
		},
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func intPtr(v int) *int         { return &v }
func f64Ptr(v float64) *float64 { return &v }
func strPtr(v string) *string   { return &v }
func timePtr(v time.Time) *time.Time { return &v }

func labelMeta(label string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"label": label})
	return b
}

func sessionMeta(session string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"session": session})
	return b
}

func sessionLabelMeta(session, label string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"session": session, "label": label})
	return b
}
