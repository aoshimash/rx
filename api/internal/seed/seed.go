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
func Run(ctx context.Context, programTemplateRepo repository.ProgramTemplateRepository, programRepo repository.ProgramRepository, logRepo repository.LogRepository) error {
	// Create SBD ProgramTemplate
	sbdTmpl := sbdProgramTemplate()
	if err := programTemplateRepo.Create(ctx, sbdTmpl); err != nil {
		return fmt.Errorf("create SBD program template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", sbdTmpl.Name)

	// Create Accessory ProgramTemplate
	accTmpl := accessoryProgramTemplate()
	if err := programTemplateRepo.Create(ctx, accTmpl); err != nil {
		return fmt.Errorf("create accessory program template: %w", err)
	}
	slog.Info("[seed] Created program template", "name", accTmpl.Name)

	// Generate SBD Program from template
	sbdInput := &domain.GenerateProgramInput{
		TargetWeights: map[string]float64{
			"スクワット":  145.0,
			"ベンチプレス": 110.0,
			"デッドリフト": 180.0,
		},
		LoadIncrements: map[string]float64{
			"スクワット":  2.5,
			"ベンチプレス": 2.5,
			"デッドリフト": 2.5,
		},
	}
	sbdProgram := domain.GenerateProgramFromTemplate(sbdTmpl, sbdInput)
	if err := programRepo.Create(ctx, sbdProgram); err != nil {
		return fmt.Errorf("create SBD program: %w", err)
	}
	slog.Info("[seed] Created program", "name", sbdProgram.Name)

	// Create Accessory Program directly (not from template)
	accProgram := accessoryProgram(accTmpl)
	if err := programRepo.Create(ctx, accProgram); err != nil {
		return fmt.Errorf("create accessory program: %w", err)
	}
	slog.Info("[seed] Created program", "name", accProgram.Name)

	// Create training logs for the first 3 SBD sessions (W1D1, W1D2, W2D1)
	// Remaining sessions (W2D2, W3D1, W3D2) have no logs — program is still active
	sessionLogs := sbdTrainingLogs(sbdProgram)
	for _, l := range sessionLogs {
		lCopy := l
		if err := logRepo.Create(ctx, &lCopy); err != nil {
			return fmt.Errorf("create log %s: %w", l.PerformedAt.Format("2006-01-02"), err)
		}
		slog.Info("[seed] Created log", "date", l.PerformedAt.Format("2006-01-02"))
	}

	return nil
}

func sbdProgramTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "SBD 3-Week プログラム",
		Description: strPtr("3週間のスクワット・ベンチプレス・デッドリフト周期プログラム（週2日）"),
		Entries: []domain.ProgramTemplateEntry{
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

func accessoryProgramTemplate() *domain.ProgramTemplate {
	return &domain.ProgramTemplate{
		Name:        "補助種目 2-Week プログラム",
		Description: strPtr("SBD の弱点補強のための補助種目（2週間サイクル・週2日）"),
		Entries: []domain.ProgramTemplateEntry{
			// ── Week 1 Day 1: Posterior Chain ──
			{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 2, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			{Order: 3, ExerciseName: "ハムストリングカール", Sets: intPtr(3), Reps: intPtr(12), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 1")},
			// ── Week 1 Day 2: Upper & Quad ──
			{Order: 4, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(3), Reps: intPtr(8), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
			{Order: 5, ExerciseName: "フロントスクワット", Sets: intPtr(3), Reps: intPtr(5), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
			{Order: 6, ExerciseName: "ダンベルショルダープレス", Sets: intPtr(3), Reps: intPtr(10), RPE: intPtr(7), Metadata: sessionMeta("Week 1 Day 2")},
		},
	}
}

// accessoryProgram creates a Program directly (not via template generation).
func accessoryProgram(tmpl *domain.ProgramTemplate) *domain.Program {
	tid := tmpl.ID
	return &domain.Program{
		ProgramTemplateID: &tid,
		Name:              "補助種目 Week 1",
		Status:            domain.ProgramStatusActive,
		Sessions: []domain.ProgramSession{
			{
				SessionName: "Week 1 Day 1",
				Order:       0,
				Entries: []domain.ProgramSessionEntry{
					{Order: 1, ExerciseName: "ルーマニアンデッドリフト", Sets: intPtr(3), Reps: intPtr(8), LoadKg: f64Ptr(80), RPE: intPtr(7)},
					{Order: 2, ExerciseName: "バーベルロウ", Sets: intPtr(4), Reps: intPtr(8), LoadKg: f64Ptr(70), RPE: intPtr(7)},
					{Order: 3, ExerciseName: "ハムストリングカール", Sets: intPtr(3), Reps: intPtr(12), LoadKg: f64Ptr(30), RPE: intPtr(7)},
				},
			},
			{
				SessionName: "Week 1 Day 2",
				Order:       1,
				Entries: []domain.ProgramSessionEntry{
					{Order: 1, ExerciseName: "クローズグリップベンチプレス", Sets: intPtr(3), Reps: intPtr(8), LoadKg: f64Ptr(80), RPE: intPtr(7)},
					{Order: 2, ExerciseName: "フロントスクワット", Sets: intPtr(3), Reps: intPtr(5), LoadKg: f64Ptr(70), RPE: intPtr(7)},
					{Order: 3, ExerciseName: "ダンベルショルダープレス", Sets: intPtr(3), Reps: intPtr(10), LoadKg: f64Ptr(20), RPE: intPtr(7)},
				},
			},
		},
	}
}

// sbdTrainingLogs creates logs for the first 3 sessions of the SBD program (W1D1, W1D2, W2D1).
func sbdTrainingLogs(program *domain.Program) []domain.Log {
	// Find session IDs by name (assigned by Create in the store)
	sessionMap := make(map[string]string) // session_name → session_name (for referencing)
	for _, s := range program.Sessions {
		sessionMap[s.SessionName] = s.SessionName
	}

	programID := program.ID
	w1d1 := "Week 1 Day 1"
	w1d2 := "Week 1 Day 2"
	w2d1 := "Week 2 Day 1"

	return []domain.Log{
		// Week 1 Day 1: Heavy (SQ 1x1/3x3/3x5, BP 1x1/3x3/3x5, DL 1x1/2x3)
		{
			ProgramID:   &programID,
			SessionName: &w1d1,
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
			ProgramID:   &programID,
			SessionName: &w1d2,
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
			ProgramID:   &programID,
			SessionName: &w2d1,
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
