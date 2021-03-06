package ts

import (
	"context"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	_ "github.com/gpmidi/TapeStats/ts/migrations"
)

func (ts *TapeStatsApp) MigrationsRun(args ...string) error {
	l := ts.Log.With().Strs("args", args).Logger()

	if len(args) == 0 {
		l.Info().Msg("No args given - Running init+up")
		if err := ts.MigrationsRun("init"); err != nil {
			return err
		}
		if err := ts.MigrationsRun("up"); err != nil {
			return err
		}
		l.Info().Msg("No args given - Done running init+up")
		return nil
	}

	l.Info().Msg("Starting Migration")
	if err := migrations.DefaultCollection.DiscoverSQLMigrations("migrations"); err != nil {
		l.Error().Err(err).Msg("Failed to read/discover SQL migrations from FS")
		return err
	}

	ctx := context.Background()
	err := ts.DB.RunInTransaction(ctx, func(tx *pg.Tx) (err error) {
		oldVersion, newVersion, err := migrations.Run(tx, args...)
		l = l.With().Int64("version.old", oldVersion).Int64("version.new", newVersion).Logger()
		l.Info().Msg("Ending Migration")
		if err != nil {
			l.Error().Err(err).Msg("Failed Migration")
			return err
		}
		return nil
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed")
		return err
	}

	l.Info().Msg("Migration successful")
	return nil
}
