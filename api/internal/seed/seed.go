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

// Run inserts sample powerlifting data into the provided repositories.
func Run(ctx context.Context, programRepo repository.ProgramRepository, planRepo repository.PlanRepository, logRepo repository.LogRepository) error {
	prog1, err := createProgram(ctx, programRepo, sbdProgram())
	if err != nil {
		return fmt.Errorf("create SBD program: %w", err)
	}
	slog.Info("[seed] Created program", "name", prog1.Name)

	prog2, err := createProgram(ctx, programRepo, accessoryProgram())
	if err != nil {
		return fmt.Errorf("create accessory program: %w", err)
	}
	slog.Info("[seed] Created program", "name", prog2.Name)

	// Create SBD plans: [0]=W1D1, [1]=W1D2, [2]=W2D1, [3]=W2D2, [4]=W3D1, [5]=W3D2
	allSbdPlans := sbdPlans(prog1.ID)
	created := make([]*domain.Plan, len(allSbdPlans))
	for i, p := range allSbdPlans {
		cp, err := createPlan(ctx, planRepo, p)
		if err != nil {
			return fmt.Errorf("create SBD plan %s: %w", p.Name, err)
		}
		slog.Info("[seed] Created plan", "name", cp.Name)
		created[i] = cp
	}

	accPlan, err := createPlan(ctx, planRepo, accessoryWeek1Plan(prog2.ID))
	if err != nil {
		return fmt.Errorf("create accessory plan: %w", err)
	}
	slog.Info("[seed] Created plan", "name", accPlan.Name)

	// Completed in order: W1D1(idx 0) → W1D2(idx 1) → W2D1(idx 2)
	// Remaining: W2D2(idx 3), W3D1(idx 4), W3D2(idx 5) — shown on Plans page
	for _, l := range trainingLogs(created[0].ID, created[1].ID, created[2].ID) {
		lCopy := l
		if err := logRepo.Create(ctx, &lCopy); err != nil {
			return fmt.Errorf("create log %s: %w", l.PerformedAt.Format("2006-01-02"), err)
		}
		slog.Info("[seed] Created log", "date", l.PerformedAt.Format("2006-01-02"))
	}

	return nil
}

func createProgram(ctx context.Context, repo repository.ProgramRepository, prog domain.Program) (*domain.Program, error) {
	if err := repo.Create(ctx, &prog); err != nil {
		return nil, err
	}
	return &prog, nil
}

func createPlan(ctx context.Context, repo repository.PlanRepository, plan domain.Plan) (*domain.Plan, error) {
	if err := repo.Create(ctx, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func sbdProgram() domain.Program {
	return domain.Program{
		ID:          uuid.New(),
		Name:        "SBD 3-Week プログラム",
		Description: strPtr("3週間のスクワット・ベンチプレス・デッドリフト周期プログラム（週2日）"),
		Entries: []domain.ProgramEntry{
			// ── Week 1 Day 1: Heavy ──
			{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
			{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
			{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
			{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			// ── Week 1 Day 2: Volume ──
			{Order: 9, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			{Order: 10, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			{Order: 11, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			// ── Week 2 Day 1: Heavy ──
			{Order: 12, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 13, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			{Order: 14, ExerciseName: "スクワット", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
			{Order: 15, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 16, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			{Order: 17, ExerciseName: "ベンチプレス", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
			{Order: 18, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
			{Order: 19, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			// ── Week 2 Day 2: Volume ──
			{Order: 20, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			{Order: 21, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			{Order: 22, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			// ── Week 3 Day 1: Peak ──
			{Order: 23, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 24, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			{Order: 25, ExerciseName: "スクワット", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
			{Order: 26, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 27, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			{Order: 28, ExerciseName: "ベンチプレス", Sets: intPtr(2), Reps: intPtr(4), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
			{Order: 29, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
			{Order: 30, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(2), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			// ── Week 3 Day 2: Deload ──
			{Order: 31, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
			{Order: 32, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
			{Order: 33, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(5), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
		},
	}
}

func accessoryProgram() domain.Program {
	return domain.Program{
		ID:          uuid.New(),
		Name:        "補助種目 2-Week プログラム",
		Description: strPtr("SBD の弱点補強のための補助種目（2週間サイクル・週2日）"),
		Entries: []domain.ProgramEntry{
			// ── Week 1 Day 1: Posterior Chain ──
			{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 2, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 3, ExerciseName: "ハムストリングカール", Sets: intPtr(3), Reps: intPtr(12), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			// ── Week 1 Day 2: Upper & Quad ──
			{Order: 4, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
			{Order: 5, ExerciseName: "フロントスクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
			{Order: 6, ExerciseName: "ダンベルショルダープレス", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
			// ── Week 2 Day 1: Posterior Chain (heavier) ──
			{Order: 7, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 1")},
			{Order: 8, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 1")},
			{Order: 9, ExerciseName: "ハムストリングカール", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 1")},
			// ── Week 2 Day 2: Upper & Quad (heavier) ──
			{Order: 10, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(4), Reps: intPtr(6), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 2")},
			{Order: 11, ExerciseName: "フロントスクワット", Sets: intPtr(4), Reps: intPtr(4), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 2")},
			{Order: 12, ExerciseName: "ダンベルショルダープレス", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(8), Metadata: sessionMeta("Week 2 Day 2")},
		},
	}
}

func conversionMeta(targetWeights map[string]float64) json.RawMessage {
	meta := map[string]interface{}{
		"conversion": map[string]interface{}{
			"target_weights": targetWeights,
		},
	}
	b, _ := json.Marshal(meta)
	return b
}

var sbdTargetWeights = map[string]float64{
	"スクワット":  145.0,
	"ベンチプレス": 110.0,
	"デッドリフト": 180.0,
}

// sbdPlans returns 6 plans in program order: W1D1, W1D2, W2D1, W2D2, W3D1, W3D2.
// W1D1/W1D2/W2D1 will have logs (completed); W2D2/W3D1/W3D2 remain unexecuted.
func sbdPlans(programID uuid.UUID) []domain.Plan {
	planMeta := conversionMeta(sbdTargetWeights)
	return []domain.Plan{
		// [0] Week 1 Day 1: Heavy (completed — has log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 1 Day 1",
			Date:        datePtr(2026, 2, 26),
			SessionName: strPtr("Week 1 Day 1"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(130), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(110), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(100), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(100), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(85), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(77.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(160), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(140), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			},
		},
		// [1] Week 1 Day 2: Volume (completed — has log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 1 Day 2",
			Date:        datePtr(2026, 2, 28),
			SessionName: strPtr("Week 1 Day 2"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(100), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
				{Order: 2, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(75), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
				{Order: 3, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(130), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			},
		},
		// [2] Week 2 Day 1: Heavy (completed — has log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 2 Day 1",
			Date:        datePtr(2026, 3, 5),
			SessionName: strPtr("Week 2 Day 1"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(132.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(112.5), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(100), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(102.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(87.5), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(77.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(162.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(142.5), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			},
		},
		// [3] Week 2 Day 2: Volume (NEXT — no log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 2 Day 2",
			Date:        datePtr(2026, 3, 7),
			SessionName: strPtr("Week 2 Day 2"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(5), LoadKg: f64Ptr(102.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
				{Order: 2, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(5), LoadKg: f64Ptr(77.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
				{Order: 3, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(4), LoadKg: f64Ptr(135), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 2", "main")},
			},
		},
		// [4] Week 3 Day 1: Peak (no log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 3 Day 1",
			Date:        datePtr(2026, 3, 12),
			SessionName: strPtr("Week 3 Day 1"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(135), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(2), LoadKg: f64Ptr(120), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(2), Reps: intPtr(4), LoadKg: f64Ptr(105), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(105), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(2), LoadKg: f64Ptr(92.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(2), Reps: intPtr(4), LoadKg: f64Ptr(80), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 3 Day 1", "backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(165), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(2), LoadKg: f64Ptr(150), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 3 Day 1", "main")},
			},
		},
		// [5] Week 3 Day 2: Deload (no log)
		{
			ID:          uuid.New(),
			ProgramID:   &programID,
			Name:        "Week 3 Day 2",
			Date:        datePtr(2026, 3, 14),
			SessionName: strPtr("Week 3 Day 2"),
			Metadata:    planMeta,
			Entries: []domain.PlanEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(90), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
				{Order: 2, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(70), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
				{Order: 3, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(120), RPE: intPtr(6), Metadata: sessionLabelMeta("Week 3 Day 2", "main")},
			},
		},
	}
}

func accessoryWeek1Plan(programID uuid.UUID) domain.Plan {
	return domain.Plan{
		ID:          uuid.New(),
		ProgramID:   &programID,
		Name:        "Week 1 Day 1",
		Date:        datePtr(2026, 3, 17),
		SessionName: strPtr("Week 1 Day 1"),
		Metadata: conversionMeta(map[string]float64{
			"ルーマニアンデッドリフト": 80.0,
			"バーベルロウ":       70.0,
			"ハムストリングカール":   30.0,
		}),
		Entries: []domain.PlanEntry{
			{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), LoadKg: f64Ptr(80), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 2, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(70), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 3, ExerciseName: "ハムストリングカール", Sets: intPtr(3), Reps: intPtr(12), LoadKg: f64Ptr(30), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
		},
	}
}

// trainingLogs creates one log per plan (1:1). Completed in sequential order.
func trainingLogs(w1d1PlanID, w1d2PlanID, w2d1PlanID uuid.UUID) []domain.Log {
	return []domain.Log{
		// Week 1 Day 1: Heavy (SQ 1x1/3x3/3x5, BP 1x1/3x3/3x5, DL 1x1/2x3)
		{
			ID:          uuid.New(),
			PlanID:      &w1d1PlanID,
			PerformedAt: time.Date(2026, 2, 26, 18, 0, 0, 0, time.Local),
			Notes:       strPtr("Week 1 SBD — 調子良い"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(122.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(105), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(95), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(92.5), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(80), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(72.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 1", "backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(150), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 1 Day 1", "top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(130), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 1 Day 1", "main")},
			},
		},
		// Week 1 Day 2: Volume (SQ 4x6, BP 4x6, DL 3x5)
		{
			ID:          uuid.New(),
			PlanID:      &w1d2PlanID,
			PerformedAt: time.Date(2026, 2, 28, 18, 0, 0, 0, time.Local),
			Notes:       strPtr("Week 1 Volume — 軽めで丁寧に"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(97.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
				{Order: 2, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(6), LoadKg: f64Ptr(72.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
				{Order: 3, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(125), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 1 Day 2", "main")},
			},
		},
		// Week 2 Day 1: Heavy (SQ 1x1/4x3/2x5, BP 1x1/4x3/2x5, DL 1x1/3x3)
		{
			ID:          uuid.New(),
			PlanID:      &w2d1PlanID,
			PerformedAt: time.Date(2026, 3, 5, 18, 0, 0, 0, time.Local),
			Notes:       strPtr("Week 2 SBD — ベンチ好調"),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(125), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(107.5), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(97.5), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(95), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(4), Reps: intPtr(3), LoadKg: f64Ptr(82.5), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(2), Reps: intPtr(5), LoadKg: f64Ptr(75), RPE: intPtr(7), Metadata: sessionLabelMeta("Week 2 Day 1", "backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(155), RPE: intPtr(9), Metadata: sessionLabelMeta("Week 2 Day 1", "top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(135), RPE: intPtr(8), Metadata: sessionLabelMeta("Week 2 Day 1", "main")},
			},
		},
	}
}

func intPtr(v int) *int         { return &v }
func f64Ptr(v float64) *float64 { return &v }
func strPtr(v string) *string   { return &v }

func datePtr(year, month, day int) *domain.DateOnly {
	d := domain.DateOnly(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
	return &d
}

func sessionLabelMeta(session, label string) json.RawMessage {
	m := map[string]string{"session": session}
	if label != "" {
		m["label"] = label
	}
	b, _ := json.Marshal(m)
	return b
}

func sessionMeta(session string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"session": session})
	return b
}
