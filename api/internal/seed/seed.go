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

	plan1, err := createPlan(ctx, planRepo, sbdWeek1Plan(prog1.ID))
	if err != nil {
		return fmt.Errorf("create SBD plan: %w", err)
	}
	slog.Info("[seed] Created plan", "name", plan1.Name)

	_, err = createPlan(ctx, planRepo, accessoryWeek1Plan(prog2.ID))
	if err != nil {
		return fmt.Errorf("create accessory plan: %w", err)
	}
	slog.Info("[seed] Created plan", "name", "補助 Week 1")

	for _, l := range trainingLogs(plan1.ID) {
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
		Name:        "SBD 基本プログラム",
		Description: strPtr("スクワット・ベンチプレス・デッドリフトを中心としたパワーリフティング向けプログラム"),
		Entries: []domain.ProgramEntry{
			// スクワット
			{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: setTypeMeta("top")},
			{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: setTypeMeta("main")},
			{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: setTypeMeta("backoff")},
			// ベンチプレス
			{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: setTypeMeta("top")},
			{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), RPE: intPtr(8), Metadata: setTypeMeta("main")},
			{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: setTypeMeta("backoff")},
			// デッドリフト
			{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), RPE: intPtr(9), Metadata: setTypeMeta("top")},
			{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), RPE: intPtr(8), Metadata: setTypeMeta("main")},
		},
	}
}

func accessoryProgram() domain.Program {
	return domain.Program{
		ID:          uuid.New(),
		Name:        "補助種目プログラム",
		Description: strPtr("SBD の弱点補強のための補助種目"),
		Entries: []domain.ProgramEntry{
			{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7)},
			{Order: 2, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7)},
			{Order: 3, ExerciseName: "フロントスクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7)},
			{Order: 4, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7)},
		},
	}
}

func sbdWeek1Plan(programID uuid.UUID) domain.Plan {
	return domain.Plan{
		ID:        uuid.New(),
		ProgramID: &programID,
		Name:      "SBD Week 1",
		Notes:     strPtr("1週目: スクワットTop130kg / ベンチTop100kg / デッドリフトTop160kg"),
		Entries: []domain.PlanEntry{
			// Day 1
			{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(130), RPE: intPtr(9), Metadata: entryMeta("Day 1", "top"), Date: datePtr(2026, 3, 16)},
			{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(110), RPE: intPtr(8), Metadata: entryMeta("Day 1", "main"), Date: datePtr(2026, 3, 16)},
			{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(100), RPE: intPtr(7), Metadata: entryMeta("Day 1", "backoff"), Date: datePtr(2026, 3, 16)},
			{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(100), RPE: intPtr(9), Metadata: entryMeta("Day 1", "top"), Date: datePtr(2026, 3, 16)},
			{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(85), RPE: intPtr(8), Metadata: entryMeta("Day 1", "main"), Date: datePtr(2026, 3, 16)},
			{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(77.5), RPE: intPtr(7), Metadata: entryMeta("Day 1", "backoff"), Date: datePtr(2026, 3, 16)},
			{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(160), RPE: intPtr(9), Metadata: entryMeta("Day 1", "top"), Date: datePtr(2026, 3, 16)},
			{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(140), RPE: intPtr(8), Metadata: entryMeta("Day 1", "main"), Date: datePtr(2026, 3, 16)},
			// Day 2
			{Order: 9, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(132.5), RPE: intPtr(9), Metadata: entryMeta("Day 2", "top"), Date: datePtr(2026, 3, 19)},
			{Order: 10, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(112.5), RPE: intPtr(8), Metadata: entryMeta("Day 2", "main"), Date: datePtr(2026, 3, 19)},
			{Order: 11, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(102.5), RPE: intPtr(7), Metadata: entryMeta("Day 2", "backoff"), Date: datePtr(2026, 3, 19)},
			{Order: 12, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(102.5), RPE: intPtr(9), Metadata: entryMeta("Day 2", "top"), Date: datePtr(2026, 3, 19)},
			{Order: 13, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(87.5), RPE: intPtr(8), Metadata: entryMeta("Day 2", "main"), Date: datePtr(2026, 3, 19)},
			{Order: 14, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(80), RPE: intPtr(7), Metadata: entryMeta("Day 2", "backoff"), Date: datePtr(2026, 3, 19)},
			{Order: 15, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(162.5), RPE: intPtr(9), Metadata: entryMeta("Day 2", "top"), Date: datePtr(2026, 3, 19)},
			{Order: 16, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(142.5), RPE: intPtr(8), Metadata: entryMeta("Day 2", "main"), Date: datePtr(2026, 3, 19)},
		},
	}
}

func accessoryWeek1Plan(programID uuid.UUID) domain.Plan {
	return domain.Plan{
		ID:        uuid.New(),
		ProgramID: &programID,
		Name:      "補助 Week 1",
		Entries: []domain.PlanEntry{
			{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), LoadKg: f64Ptr(80), RPE: intPtr(7), Metadata: sessionMeta("補助 Day 1"), Date: datePtr(2026, 3, 17)},
			{Order: 2, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(3), Reps: intPtr(8), LoadKg: f64Ptr(65), RPE: intPtr(7), Metadata: sessionMeta("補助 Day 1"), Date: datePtr(2026, 3, 17)},
			{Order: 3, ExerciseName: "フロントスクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(80), RPE: intPtr(7), Metadata: sessionMeta("補助 Day 1"), Date: datePtr(2026, 3, 17)},
			{Order: 4, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(70), RPE: intPtr(7), Metadata: sessionMeta("補助 Day 1"), Date: datePtr(2026, 3, 17)},
		},
	}
}

func trainingLogs(planID uuid.UUID) []domain.Log {
	type weekData struct {
		date      time.Time
		squatTop  float64
		squatMain float64
		squatBO   float64
		benchTop  float64
		benchMain float64
		benchBO   float64
		deadTop   float64
		deadMain  float64
	}
	weeks := []weekData{
		{time.Date(2026, 2, 26, 18, 0, 0, 0, time.Local), 122.5, 105, 95, 92.5, 80, 72.5, 150, 130},
		{time.Date(2026, 3, 5, 18, 0, 0, 0, time.Local), 125, 107.5, 97.5, 95, 82.5, 75, 155, 135},
		{time.Date(2026, 3, 12, 18, 0, 0, 0, time.Local), 127.5, 110, 100, 97.5, 85, 77.5, 157.5, 137.5},
	}
	logs := make([]domain.Log, 0, len(weeks))
	for i, w := range weeks {
		planIDCopy := planID
		logs = append(logs, domain.Log{
			ID:          uuid.New(),
			PlanID:      &planIDCopy,
			PerformedAt: w.date,
			Notes:       strPtr(fmt.Sprintf("Week %d SBD", i+1)),
			Entries: []domain.LogEntry{
				{Order: 1, ExerciseName: "スクワット", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(w.squatTop), RPE: intPtr(9), Metadata: setTypeMeta("top")},
				{Order: 2, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(w.squatMain), RPE: intPtr(8), Metadata: setTypeMeta("main")},
				{Order: 3, ExerciseName: "スクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(w.squatBO), RPE: intPtr(7), Metadata: setTypeMeta("backoff")},
				{Order: 4, ExerciseName: "ベンチプレス", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(w.benchTop), RPE: intPtr(9), Metadata: setTypeMeta("top")},
				{Order: 5, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(3), LoadKg: f64Ptr(w.benchMain), RPE: intPtr(8), Metadata: setTypeMeta("main")},
				{Order: 6, ExerciseName: "ベンチプレス", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(w.benchBO), RPE: intPtr(7), Metadata: setTypeMeta("backoff")},
				{Order: 7, ExerciseName: "デッドリフト", Sets: intPtr(1), Reps: intPtr(1), LoadKg: f64Ptr(w.deadTop), RPE: intPtr(9), Metadata: setTypeMeta("top")},
				{Order: 8, ExerciseName: "デッドリフト", Sets: intPtr(2), Reps: intPtr(3), LoadKg: f64Ptr(w.deadMain), RPE: intPtr(8), Metadata: setTypeMeta("main")},
			},
		})
	}
	return logs
}

func intPtr(v int) *int         { return &v }
func f64Ptr(v float64) *float64 { return &v }
func strPtr(v string) *string   { return &v }

func datePtr(year, month, day int) *domain.DateOnly {
	d := domain.DateOnly(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
	return &d
}

func sessionMeta(session string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"session": session})
	return b
}

func setTypeMeta(setType string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"set_type": setType})
	return b
}

func entryMeta(session, setType string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{"session": session, "set_type": setType})
	return b
}
