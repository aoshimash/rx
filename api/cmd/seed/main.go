package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/aoshimash/rx/api/internal/config"
	"github.com/aoshimash/rx/api/internal/seed"
	"github.com/aoshimash/rx/api/internal/store/memory"
	postgresstore "github.com/aoshimash/rx/api/internal/store/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	reset := flag.Bool("reset", false, "truncate all tables before seeding (postgres only)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()
	cfg := config.Load()

	if cfg.Database.StorageType == "postgres" {
		db, err := postgresstore.NewDB(ctx, cfg.Database)
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		if *reset {
			slog.Warn("Resetting database (TRUNCATE CASCADE)...")
			if err := resetDB(ctx, db.Pool()); err != nil {
				slog.Error("Failed to reset database", "error", err)
				os.Exit(1)
			}
			slog.Info("Database reset complete")
		}

		pool := db.Pool()
		err = seed.Run(ctx,
			postgresstore.NewProgramTemplateRepository(pool),
			postgresstore.NewProgramRepository(pool),
			postgresstore.NewLogRepository(pool),
		)
		if err != nil {
			slog.Error("Seeding failed", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Warn("STORAGE_TYPE is not postgres — seeding in-memory (data will not persist)")
		err := seed.Run(ctx,
			memory.NewProgramTemplateRepository(),
			memory.NewProgramRepository(),
			memory.NewLogRepository(),
		)
		if err != nil {
			slog.Error("Seeding failed", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("Seeding completed successfully")
}

func resetDB(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx,
		"TRUNCATE log_entries, logs, program_session_entries, program_sessions, programs, program_template_entries, program_templates RESTART IDENTITY CASCADE",
	)
	return err
}
